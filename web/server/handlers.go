package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func root(c *gin.Context) {
	user_id := c.GetString("user_id")
	if user_id != "" {
		c.Redirect(http.StatusSeeOther, "/quotes")
	} else {
		c.HTML(http.StatusOK, "auth.html", nil)
	}
}

func quotes(c *gin.Context) {
	c.HTML(http.StatusOK, "quotes.html", nil)
}

func settings(c *gin.Context) {
	c.HTML(http.StatusOK, "settings.html", nil)
}
