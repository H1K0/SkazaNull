package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func root(c *gin.Context) {
	authorized := c.GetBool("authorized")
	if authorized {
		c.Redirect(http.StatusSeeOther, "/quotes")
	} else {
		c.HTML(http.StatusOK, "auth.html", nil)
	}
}

func quotes(c *gin.Context) {
	authorized := c.GetBool("authorized")
	if authorized {
		c.HTML(http.StatusOK, "quotes.html", nil)
	} else {
		c.Redirect(http.StatusSeeOther, "/")
	}
}
