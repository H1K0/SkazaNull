package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

func HandleDBError(err error) (int, string) {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return 500, err.Error()
	}
	statusStr := pgErr.Message[:3]
	message := pgErr.Message[4:]
	status, err := strconv.ParseInt(statusStr, 10, 0)
	if err == nil {
		return int(status), message
	} else {
		return 400, pgErr.Message
	}
}

func MiddlewareAuth(c *gin.Context) {
	session := sessions.Default(c)
	user_id, ok := session.Get("user_id").(string)
	if !ok && strings.HasPrefix(c.Request.URL.Path, "/api") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Ты это, залогинься сначала что ли, а то чё как крыса"})
		return
	}
	c.Set("authorized", ok)
	c.Set("user_id", user_id)
	c.Next()
}
