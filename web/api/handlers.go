package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/H1K0/SkazaNull/db"
	"github.com/H1K0/SkazaNull/models"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "И чё за шнягу ты мне кинул?"})
		return
	}
	user, err := db.UserAuth(context.Background(), credentials.Login, credentials.Password)
	if err != nil {
		status, message := handleDBError(err)
		c.JSON(status, gin.H{"error": message})
		return
	}
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("started", time.Now().Unix())
	session.Save()
	c.JSON(http.StatusOK, user)
}

func userGet(c *gin.Context) {
	user_id := c.GetString("user_id")
	user, err := db.UserGet(context.Background(), user_id)
	if err != nil {
		status, message := handleDBError(err)
		c.JSON(status, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, user)
}

func userUpdate(c *gin.Context) {
	user_id := c.GetString("user_id")
	var body map[string]string
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "И чё за шнягу ты мне кинул?"})
		return
	}
	var user models.User
	ctx := context.Background()
	newTelegramIDStr, ok := body["telegram_id"]
	if ok && newTelegramIDStr != "" {
		newTelegramID, err := strconv.ParseInt(newTelegramIDStr, 10, 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "telegram_id param must be int"})
			return
		}
		user, err = db.UserUpdateTelegramID(ctx, user_id, newTelegramID)
		if err != nil {
			status, message := handleDBError(err)
			c.JSON(status, gin.H{"error": message})
			return
		}
	}
	newName, ok := body["name"]
	if ok && newName != "" {
		user, err = db.UserUpdateName(ctx, user_id, newName)
		if err != nil {
			status, message := handleDBError(err)
			c.JSON(status, gin.H{"error": message})
			return
		}
	}
	newLogin, ok := body["login"]
	if ok && newLogin != "" {
		user, err = db.UserUpdateLogin(ctx, user_id, newLogin)
		if err != nil {
			status, message := handleDBError(err)
			c.JSON(status, gin.H{"error": message})
			return
		}
	}
	newPassword, ok := body["password"]
	if ok && newPassword != "" {
		user, err = db.UserUpdatePassword(ctx, user_id, newPassword)
		if err != nil {
			status, message := handleDBError(err)
			c.JSON(status, gin.H{"error": message})
			return
		}
	}
	c.JSON(http.StatusOK, user)
}

func userLogout(c *gin.Context) {
	session := sessions.Default(c)
	_, ok := session.Get("user_id").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ну и как я тебя разлогиню, если ты даже не залогинился?"})
		return
	}
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()
	c.Status(http.StatusNoContent)
}

//#endregion User

//#region Quotes

func quotesGet(c *gin.Context) {
	user_id := c.GetString("user_id")
	filter, ok := c.GetQuery("filter")
	if !ok {
		filter = ""
	}
	sort, ok := c.GetQuery("sort")
	if !ok {
		sort = "-datetime"
	}
	var limit int64
	limitStr, ok := c.GetQuery("limit")
	if ok {
		var err error
		limit, err = strconv.ParseInt(limitStr, 10, 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit param must be int"})
			return
		}
	} else {
		limit = 10
	}
	var offset int64
	offsetStr, ok := c.GetQuery("offset")
	if ok {
		var err error
		offset, err = strconv.ParseInt(offsetStr, 10, 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "offset param must be int"})
			return
		}
	} else {
		offset = 0
	}
	quotes, err := db.QuotesGet(c, user_id, filter, sort, int(limit), int(offset))
	if err != nil {
		status, message := handleDBError(err)
		c.JSON(status, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, quotes)
}

func quoteGet(c *gin.Context) {
	user_id := c.GetString("user_id")
	quote_id := c.Param("id")
	quote, err := db.QuoteGet(context.Background(), user_id, quote_id)
	if err != nil {
		status, message := handleDBError(err)
		c.JSON(status, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, quote)
}

func quoteAdd(c *gin.Context) {
	user_id := c.GetString("user_id")
	var body map[string]string
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "И чё за шнягу ты мне кинул?"})
		return
	}
	text, ok := body["text"]
	if !ok || text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Э, а где цитата?"})
		return
	}
	author, ok := body["author"]
	if !ok || author == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Цитата может быть сказана только тем, кто её сказанул! А кто сказанул эту цитату?"})
		return
	}
	var datetime time.Time
	datetimeStr, ok := body["datetime"]
	if ok {
		datetime, err = time.Parse(time.RFC3339, datetimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Чёт дата и время у тебя какие-то кривые..."})
			return
		}
	} else {
		datetime = time.Now()
	}
	quote, err := db.QuoteAdd(context.Background(), user_id, text, author, datetime)
	if err != nil {
		status, message := handleDBError(err)
		c.JSON(status, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusCreated, quote)
}

func quoteUpdate(c *gin.Context) {
	user_id := c.GetString("user_id")
	quote_id := c.Param("id")
	var body map[string]string
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "И чё за шнягу ты мне кинул?"})
		return
	}
	var quote models.Quote
	ctx := context.Background()
	newText, ok := body["text"]
	if ok && newText != "" {
		quote, err = db.QuoteUpdateText(ctx, user_id, quote_id, newText)
		if err != nil {
			status, message := handleDBError(err)
			c.JSON(status, gin.H{"error": message})
			return
		}
	}
	newAuthor, ok := body["author"]
	if ok && newAuthor != "" {
		quote, err = db.QuoteUpdateAuthor(ctx, user_id, quote_id, newAuthor)
		if err != nil {
			status, message := handleDBError(err)
			c.JSON(status, gin.H{"error": message})
			return
		}
	}
	datetimeStr, ok := body["datetime"]
	if ok && datetimeStr != "" {
		newDatetime, err := time.Parse(time.RFC3339, datetimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Чёт дата и время у тебя какие-то кривые..."})
			return
		}
		quote, err = db.QuoteUpdateDatetime(context.Background(), user_id, quote_id, newDatetime)
		if err != nil {
			status, message := handleDBError(err)
			c.JSON(status, gin.H{"error": message})
			return
		}
	}
	c.JSON(http.StatusOK, quote)
}

func quoteDelete(c *gin.Context) {
	user_id := c.GetString("user_id")
	quote_id := c.Param("id")
	err := db.QuoteDelete(context.Background(), user_id, quote_id)
	if err != nil {
		status, message := handleDBError(err)
		c.JSON(status, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

//#endregion Quotes
