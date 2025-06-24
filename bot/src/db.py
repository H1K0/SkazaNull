from psycopg2 import connect as db_connect
from psycopg2.extras import DictCursor


class Database:
	def __init__(self, **kwargs):
		self.conn = db_connect(**kwargs)
		self.conn.autocommit = True

	def close(self):
		self.conn.close()

	# get all users
	def get_users(self):
		with self.conn.cursor(cursor_factory=DictCursor) as cursor:
			cursor.execute("SELECT * FROM users")
			return cursor.fetchall()

	# check if message sender is authorized
	def check_user(self, uid):
		with self.conn.cursor() as cursor:
			cursor.execute("SELECT 1 FROM users WHERE telegram_id=%s", (uid,))
			return cursor.fetchone() is not None

	# check if user is an editor
	def check_editor(self, msg):
		with self.conn.cursor() as cursor:
			cursor.execute("SELECT is_editor FROM users WHERE telegram_id=%s", (msg.from_user.id,))
			return cursor.fetchone()[0]

	# get quotes range
	def get_quotes(self, start=0, count="ALL"):
		with self.conn.cursor(cursor_factory=DictCursor) as cursor:
			cursor.execute("SELECT q.id, q.text, to_char(q.datetime, 'DD.MM.YYYY HH24:MI:SS') AS datetime, q.author, q.creator_id, u.name AS creator_name FROM quotes q JOIN users u ON q.creator_id = u.id WHERE NOT q.is_removed ORDER BY q.datetime DESC OFFSET %s LIMIT %s", (start, count))
			return cursor.fetchall()

	# get single quote
	def get_quote(self, id):
		with self.conn.cursor(cursor_factory=DictCursor) as cursor:
			cursor.execute("SELECT q.id, q.text, to_char(q.datetime, 'DD.MM.YYYY HH24:MI:SS') AS datetime, q.author, q.creator_id, u.name AS creator_name FROM quotes q JOIN users u ON q.creator_id = u.id WHERE NOT q.is_removed AND q.id=%s", (id,))
			return cursor.fetchone()

	# get random quote
	def get_random_quote(self):
		with self.conn.cursor(cursor_factory=DictCursor) as cursor:
			cursor.execute("SELECT q.id, q.text, to_char(q.datetime, 'DD.MM.YYYY HH24:MI:SS') AS datetime, q.author, q.creator_id, u.name AS creator_name FROM quotes q JOIN users u ON q.creator_id = u.id ORDER BY random() LIMIT 1")
			return cursor.fetchone()

	# add quote
	def add_quote(self, msg):
		if not self.check_editor(msg):
			raise RuntimeError
		quote = msg.text.strip().split("\n\n")
		text = quote[0].strip()
		author = "" if len(quote) == 1 else quote[1].strip()
		timestamp = msg.forward_origin.date if msg.forward_origin else msg.date
		with self.conn.cursor() as cursor:
			cursor.execute("INSERT INTO quotes(text, author, creator_id, datetime) VALUES(%s, NULLIF(%s, ''), (SELECT id FROM users WHERE telegram_id=%s), to_timestamp(%s)::timestamptz) RETURNING id", (quote[0].strip(), author, msg.from_user.id, timestamp))
			return cursor.fetchone()[0]

	def quotes_count(self):
		with self.conn.cursor() as cursor:
			cursor.execute("SELECT count(*) FROM quotes")
			return cursor.fetchone()[0]
