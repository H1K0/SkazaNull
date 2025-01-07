package server

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/H1K0/SkazaNull/api"
	"github.com/H1K0/SkazaNull/embed"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Serve(addr string, encryptionKey []byte) {
	r := gin.Default()

	store := cookie.NewStore(encryptionKey)
	store.Options(sessions.Options{Path: "/"})
	r.Use(sessions.Sessions("session", store))

	tmpl := template.Must(template.ParseFS(embed.TemplatesFS, "templates/*.html"))
	r.SetHTMLTemplate(tmpl)

	api.RegisterRoutes(r)

	static, err := fs.Sub(embed.StaticFS, "static")
	if err != nil {
		log.Fatalf("Failed to get subs of embedded static FS: %s\n", err)
	}
	r.StaticFS("/static/", http.FS(static))
	r.StaticFileFS("/favicon.ico", "static/service/favicon.ico", http.FS(embed.StaticFS))
	r.StaticFileFS("/skazanull.webmanifest", "static/service/skazanull.webmanifest", http.FS(embed.StaticFS))
	r.StaticFileFS("/browserconfig.xml", "static/service/browserconfig.xml", http.FS(embed.StaticFS))

	r.GET("/", api.MiddlewareAuth, root)
	r.GET("/quotes", api.MiddlewareAuth, middlewareAuth, quotes)
	r.GET("/settings", api.MiddlewareAuth, middlewareAuth, settings)

	r.Run(addr)
}
