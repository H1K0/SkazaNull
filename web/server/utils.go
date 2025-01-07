package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func middlewareAuth(c *gin.Context) {
	user_id := c.GetString("user_id")
	if user_id == "" {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	c.Next()
}
