package api

import (
	"github.com/gin-gonic/gin"
)

// @title SkazaNull API
// @description RESTful API для пацанского цитатника SkazaNull
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host skazanull.hakoniwa.ru
// @BasePath /api
func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/auth", userAuth)
		api.GET("/auth", MiddlewareAuth, userGet)
		api.PATCH("/auth", MiddlewareAuth, userUpdate)
		api.DELETE("/auth", userLogout)

		api.GET("/quotes", MiddlewareAuth, quotesGet)
		api.POST("/quotes", MiddlewareAuth, quoteAdd)
		api.GET("/quotes/:id", MiddlewareAuth, quoteGet)
		api.PATCH("/quotes/:id", MiddlewareAuth, quoteUpdate)
		api.DELETE("/quotes/:id", MiddlewareAuth, quoteDelete)
	}
}
