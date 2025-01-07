package server

import (
	"github.com/H1K0/SkazaNull/api"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Serve(addr string, encryptionKey []byte) {
	r := gin.Default()

	store := cookie.NewStore(encryptionKey)
	store.Options(sessions.Options{Path: "/"})
	r.Use(sessions.Sessions("session", store))

	api.RegisterRoutes(r)

	r.LoadHTMLGlob("templates/*.html")

	r.Static("/favicon.ico", "./static/service/favicon.ico")
	r.Static("/skazanull.webmanifest", "./static/service/skazanull.webmanifest")
	r.Static("/browserconfig.xml", "./static/service/browserconfig.xml")
	r.Static("/static", "./static")

	r.GET("/", api.MiddlewareAuth, root)
	r.GET("/quotes", api.MiddlewareAuth, middlewareAuth, quotes)
	r.GET("/settings", api.MiddlewareAuth, middlewareAuth, settings)

	r.Run(addr)
}
