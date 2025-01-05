from telebot import TeleBot
from telebot.types import InlineKeyboardMarkup, InlineKeyboardButton
from configparser import ConfigParser
from db import Database
import os
import atexit
import signal
import logging as log
from time import sleep

# set logger
log.basicConfig(
	level=log.INFO,
	filename="/var/log/skazanull/bot.log",
	filemode="a",
	format="%(asctime)s | %(levelname)s | %(message)s"
)

# actions to do on exit
exit_actions = []
defer = exit_actions.append
def finalize(*args):
	try:
		exec('\n'.join(exit_actions))
	finally:
		os._exit(0)
atexit.register(finalize)
signal.signal(signal.SIGTERM, finalize)
signal.signal(signal.SIGINT, finalize)

# load config
try:
	conf = ConfigParser()
	conf.read("/etc/skazanull/bot.conf")
except Exception as e:
	log.critical(f"failed to load config: {str(e)}")
	exit(1)

# connect to db
try:
	db = Database(
		host=conf["DB"]["Host"],
		port=conf["DB"]["Port"],
		dbname=conf["DB"]["Name"],
		user=conf["DB"]["User"],
		password=conf["DB"]["Password"],
		application_name="SkazaNull Telegram Bot",
	)
	defer("db.close()")
except Exception as e:
	log.critical(f"failed to connect to db: {str(e)}")
	exit(1)

# initialize bot
snb = TeleBot(conf["General"]["Token"])
memo = {}

# print help
@snb.message_handler(commands=["start", "help"])
def helper(msg):
	if not db.check_user(msg.from_user.id):
		log.info(f"unauthorized access from user {msg.from_user.id}")
		snb.send_message(msg.chat.id, "Я не понял, ты вообще кто такой?")
		return
	snb.send_message(msg.chat.id, """*SkazaNull - пацанский ботяра для пацанских цитат*

/help - Помощь нужна тому, кто ее просит, а не тому, кто ее не просит
/random - Показать рандомную цитату из базы
/quotes - Перейти в режим просмотра пацанских цитат

Чтобы добавить новую цитату, просто пришли ее текстовым сообщением. Авторов цитаты обязательно укажи в конце, отделив двумя переносами строк. Такие дела, бро.

Пример оформления цитаты:
```
- А что это за жопа?
- Так это же Гурилий!

Биба и Боба
```""", parse_mode="markdown")

# quote formatter
def format_quote(quote):
	return "Добавил %s %s:\n```\n%s\n%s```" % (quote["creator_name"], quote["datetime"], quote["text"], f"\n© {quote['author']}" if quote["author"] else "")

# quotes view
@snb.message_handler(commands=["quotes"])
def quotes_handler(msg, page=0, prev_msg=None, user=None):
	if not user:
		user = msg.from_user.id
	if not db.check_user(user):
		log.info(f"unauthorized access from user {user}")
		snb.send_message(msg.chat.id, "Я не понял, ты вообще че за кернел??")
		return
	quotes_count = db.quotes_count()
	if quotes_count == 0:
		snb.send_message(msg.chat.id, "А где цитаты-то?")
		return
	page_size = int(conf["Output"]["QuotesPerPage"])
	pages_count = quotes_count // page_size + bool(quotes_count % page_size)
	prev_page = (page-1) % pages_count
	next_page = (page+1) % pages_count
	buttons = InlineKeyboardMarkup()
	prev_button = InlineKeyboardButton("←", callback_data=f"quotes:{prev_page}")
	curr_button = InlineKeyboardButton(f"{page*page_size+1}-{min(quotes_count,(page+1)*page_size)}/{quotes_count}", callback_data=":")
	next_button = InlineKeyboardButton("→", callback_data=f"quotes:{next_page}")
	buttons.add(prev_button, curr_button, next_button)
	snb.send_message(msg.chat.id, "\n\n".join(list(map(format_quote, db.get_quotes(page_size*page, page_size)))), parse_mode="markdown", reply_markup=buttons)
	if prev_msg:
		snb.delete_message(msg.chat.id, prev_msg.id)

# search quotes by keywords
@snb.message_handler(commands=["search"])
def search_quotes(msg):
	pass

# get random quote
@snb.message_handler(commands=["random"])
def random_quote(msg):
	if not db.check_user(msg.from_user.id):
		log.info(f"unauthorized access from user {msg.from_user.id}")
		snb.send_message(msg.chat.id, "Я не понял, ты вообще кто такой?")
		return
	quote = db.get_random_quote()
	if quote:
		snb.send_message(msg.chat.id, format_quote(quote), parse_mode="markdown")
		return
	snb.send_message(msg.chat.id, "А где цитаты-то?")

# add new quote
@snb.message_handler(content_types=["text"])
def text_handler(msg):
	if not db.check_user(msg.from_user.id):
		log.info(f"unauthorized access from user {msg.from_user.id}")
		snb.send_message(msg.chat.id, "Я не понял, ты вообще че за кернел??")
		return
	try:
		qid = db.add_quote(msg)
		log.info(f"quote '{qid}' added")
		snb.reply_to(msg, "Цитата добавлена!")
		quote = db.get_quote(qid)
		for user in db.get_users():
			if user["telegram_id"] != msg.from_user.id:
				snb.send_message(user["telegram_id"], "Пользователь %s только что добавил новую цитату:\n```\n%s\n%s```" % (quote["creator_name"], quote["text"], f"\n© {quote['author']}" if quote['author'] else ""), parse_mode="markdown")
				log.info(f"notified user '{user['telegram_id']}' about quote '{qid}'")
	except RuntimeError:
		log.info(f"attempt to add quote from non-editor {msg.from_user.id}")
		snb.reply_to(msg, "Сорямбус, бро, но твои права не скачаны. Со всеми вопросами - к разрабам.")

# callbacks
@snb.callback_query_handler(lambda c: True)
def callback_handler(c):
	call = c.data.split(':')
	if call[0] == "quotes":
		quotes_handler(c.message, page=int(call[1]), prev_msg=c.message, user=c.from_user.id)


if __name__ == '__main__':
	log.info("snb started")
	defer("log.info(\"snb stopped\")")
	while 1:
		try:
			snb.polling()
		except Exception as e:
			log.error("polling stopped: " + str(e))
			sleep(1)
