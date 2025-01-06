package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/H1K0/SkazaNull/db"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//#region User

func userAuth(c *gin.Context) {
	var credentials struct {
		Login    string `form:"login"    binding:"required"`
		Password string `form:"password" binding:"required"`
	}
	err := c.ShouldBind(&credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "И че за шнягу ты мне кинул?"})
		return
	}
	user, err := db.UserAuth(context.Background(), credentials.Login, credentials.Password)
	if err != nil {
		pqErr := db.CastToPgError(err)
		if pqErr == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		statusStr := pqErr.Message[:3]
		msg := pqErr.Message[4:]
		status, _err := strconv.ParseInt(statusStr, 10, 0)
		if _err == nil {
			c.JSON(int(status), gin.H{"error": msg})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("started", time.Now().Unix())
	session.Save()
	c.JSON(http.StatusOK, user)
}

//#endregion User

//#region Quotes

//#endregion Quotes
