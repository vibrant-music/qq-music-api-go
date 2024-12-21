package routes

import (
	"github.com/gin-gonic/gin"
	"qq-music-api/controllers"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/search", controllers.Search)
	r.GET("/search/hot", controllers.HotSearch)
	r.GET("/search/quick", controllers.QuickSearch)
}
