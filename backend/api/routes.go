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
		api.GET("/quotes", quotesGet)
		api.GET("/quotes/:id", quoteGet)
	}
}
