package internalhttp

import (
	"github.com/gin-gonic/gin"
)

func NewGinRouter(app Application) *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		loggingMiddleware(c, app)
		appMiddleware(c, app)
	})


	router.POST("/banner-place", func(c *gin.Context) {
		AddBannerToPlaceHandler(c, app)
	})
	router.DELETE("/banner-place", func(c *gin.Context) {
		DeleteBannerFromPlaceHandler(c, app)
	})
	router.POST("/event", func(c *gin.Context) {
		CreateEventHandler(c, app)
	})
	router.GET("/best-banner", func(c *gin.Context) {
		BannerHandler(c, app)
	})

	return router
}
