package config

import (
	"github.com/gin-gonic/gin"
	"github.com/lucasbrito3001/url_shortner/internal/adapters/in/http"
)

func SetRoutes(router *gin.Engine, shortenedUrlController http.ShortenedUrlController) {
	v1 := router.Group("/api/v1")
	v1.POST("/shorten", shortenedUrlController.ShortenUrl)
	
	router.GET("/:code", shortenedUrlController.RedirectUrl)
}
