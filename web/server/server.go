package server

import (
	"github.com/H1K0/SkazaNull/api"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Serve(addr string) {
	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{Path: "/"})
	r.Use(sessions.Sessions("session", store))

	api.RegisterRoutes(r)

	r.LoadHTMLGlob("templates/*.html")

	r.Static("/static", "./static")
	r.GET("/", api.MiddlewareAuth, root)
	r.GET("/quotes", api.MiddlewareAuth, quotes)

	r.Run(addr)
}