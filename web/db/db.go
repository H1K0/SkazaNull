package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/H1K0/SkazaNull/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ConnPool *pgxpool.Pool

func InitDB(connString string) error {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return fmt.Errorf("error while parsing connection string: %w", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 0
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.HealthCheckPeriod = 30 * time.Second

	ConnPool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return fmt.Errorf("error while initializing DB connections pool: %w", err)
	}
	return nil
}

//#region User

func UserAuth(ctx context.Context, login string, password string) (user models.User, err error) {
	row := ConnPool.QueryRow(ctx, "SELECT * FROM user_auth($1, $2)", login, password)
	err = row.Scan(&user.ID, &user.Name, &user.Login, &user.Role, &user.TelegramID)
	return
}

func UserGet(ctx context.Context, user_id string) (user models.User, err error) {
	row := ConnPool.QueryRow(ctx, "SELECT * FROM user_get($1)", user_id)
	err = row.Scan(&user.ID, &user.Name, &user.Login, &user.Role, &user.TelegramID)
	return
}

func UserUpdateName(ctx context.Context, user_id string, new_name string) (user models.User, err error) {
	row := ConnPool.QueryRow(ctx, "SELECT * FROM user_update($1, $2)", user_id, new_name)
	err = row.Scan(&user.ID, &user.Name, &user.Login, &user.Role, &user.TelegramID)
	return
}

func UserUpdateLogin(ctx context.Context, user_id string, new_login string) (user models.User, err error) {
	row := ConnPool.QueryRow(ctx, "SELECT * FROM user_update($1, NULL, $2)", user_id, new_login)
	err = row.Scan(&user.ID, &user.Name, &user.Login, &user.Role, &user.TelegramID)
	return
}

func UserUpdateTelegramID(ctx context.Context, user_id string, new_telegram_id int64) (user models.User, err error) {
	row := ConnPool.QueryRow(ctx, "SELECT * FROM user_update($1, NULL, NULL, $2)", user_id, new_telegram_id)
	err = row.Scan(&user.ID, &user.Name, &user.Login, &user.Role, &user.TelegramID)
	return
}

func UserUpdatePassword(ctx context.Context, user_id string, new_password string) (user models.User, err error) {
	row := ConnPool.QueryRow(ctx, "SELECT * FROM user_update($1, NULL, NULL, NULL, $2)", user_id, new_password)
	err = row.Scan(&user.ID, &user.Name, &user.Login, &user.Role, &user.TelegramID)
	return
}

//#endregion User

//#region Quotes

func QuotesGet(ctx context.Context, user_id string, filter string, sort string, limit int, offset int) (quotes models.Quotes, err error) {
	queryGet := "SELECT * FROM quotes_get($1) WHERE position($2 in lower(text))>0 OR position($2 in lower(author))>0"
	if sort == "random" {
		queryGet += " ORDER BY random()"
	} else if sort != "" {
		sort_options := strings.Split(sort, ",")
		queryGet += " ORDER BY "
		for i, sort_option := range sort_options {
			sort_order := sort_option[:1]
			sort_field := sort_option[1:]
			switch sort_order {
			case "+":
				sort_order = "ASC"
			case "-":
				sort_order = "DESC"
			default:
				err = fmt.Errorf("invalid sorting parameter: %q", sort)
				return
			}
			switch sort_field {
			case "text":
				fallthrough
			case "author":
				fallthrough
			case "datetime":
				fallthrough
			case "creator.name":
				sort_field = strings.ReplaceAll(sort_field, ".", "_")
			default:
				err = fmt.Errorf("unknown sorting field: %q", sort_field)
				return
			}
			if i > 0 {
				queryGet += ", "
			}
			queryGet += fmt.Sprintf("%s %s", sort_field, sort_order)
		}
	}
	queryCount := queryGet
	if limit >= 0 {
		queryGet += fmt.Sprintf(" LIMIT %d", limit)
	}
	if offset > 0 {
		queryGet += fmt.Sprintf(" OFFSET %d", offset)
	}
	filter = strings.ToLower(filter)
	rows, err := ConnPool.Query(ctx, queryGet, user_id, filter)
	if err != nil {
		err = fmt.Errorf("error while getting quotes: %w", err)
		return
	}
	quotes.Quotes = []models.Quote{}
	count := 0
	for rows.Next() {
		var quote models.Quote
		err = rows.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Datetime, &quote.Creator.ID, &quote.Creator.Name, &quote.Creator.Login, &quote.Creator.Role, &quote.Creator.TelegramID)
		if err != nil {
			err = fmt.Errorf("error while fetching quotes: %w", err)
			return
		}
		quotes.Quotes = append(quotes.Quotes, quote)
		count++
	}
	err = rows.Err()
	if err != nil {
		return
	}
	quotes.Pagination.Limit = limit
	quotes.Pagination.Offset = offset
	quotes.Pagination.Count = count
	queryCount = fmt.Sprintf("SELECT count(*) FROM (%s) q", queryCount)
	row := ConnPool.QueryRow(ctx, queryCount, user_id, filter)
	err = row.Scan(&quotes.Pagination.TotalCount)
	return
}

func QuoteGet(ctx context.Context, user_id string, quote_id string) (quote models.Quote, err error) {
	row := ConnPool.QueryRow(ctx, "SELECT * FROM quote_get($1, $2)", user_id, quote_id)
	err = row.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Datetime, &quote.Creator.ID, &quote.Creator.Name, &quote.Creator.Login, &quote.Creator.Role, &quote.Creator.TelegramID)
	return
}

func QuoteAdd(ctx context.Context, user_id string, text string, author string, datetime time.Time) (quote models.Quote, err error) {
	row := ConnPool.QueryRow(ctx, "SELECT * FROM quote_add($1, $2, $3, $4)", user_id, text, author, datetime)
	err = row.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Datetime, &quote.Creator.ID, &quote.Creator.Name, &quote.Creator.Login, &quote.Creator.Role, &quote.Creator.TelegramID)
	return
}

func QuoteUpdateText(ctx context.Context, user_id string, quote_id string, new_text string) (quote models.Quote, err error) {
	row := ConnPool.QueryRow(ctx, "SELECT * FROM quote_update($1, $2, $3)", user_id, quote_id, new_text)
	err = row.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Datetime, &quote.Creator.ID, &quote.Creator.Name, &quote.Creator.Login, &quote.Creator.Role, &quote.Creator.TelegramID)
	return
}

func QuoteUpdateAuthor(ctx context.Context, user_id string, quote_id string, new_author string) (quote models.Quote, err error) {
	row := ConnPool.QueryRow(ctx, "SELECT * FROM quote_update($1, $2, NULL, $3)", user_id, quote_id, new_author)
	err = row.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Datetime, &quote.Creator.ID, &quote.Creator.Name, &quote.Creator.Login, &quote.Creator.Role, &quote.Creator.TelegramID)
	return
}

func QuoteUpdateDatetime(ctx context.Context, user_id string, quote_id string, new_datetime time.Time) (quote models.Quote, err error) {
	row := ConnPool.QueryRow(ctx, "SELECT * FROM quote_update($1, $2, NULL, NULL, $3)", user_id, quote_id, new_datetime)
	err = row.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Datetime, &quote.Creator.ID, &quote.Creator.Name, &quote.Creator.Login, &quote.Creator.Role, &quote.Creator.TelegramID)
	return
}

func QuoteDelete(ctx context.Context, user_id string, quote_id string) (err error) {
	_, err = ConnPool.Exec(ctx, "CALL quote_delete($1, $2)", user_id, quote_id)
	return
}

//#endregion Quotes
