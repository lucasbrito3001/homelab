package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lucasbrito3001/url_shortner/internal/app/ports"
)

type ShortenedUrlController struct {
	shortenUrlUseCase  ports.ShortenUrlUseCase
	redirectUrlUseCase ports.RedirectUrlUseCase
}

func NewShortenedUrlController(shortenUrlUseCase ports.ShortenUrlUseCase, redirectUrlUseCase ports.RedirectUrlUseCase) *ShortenedUrlController {
	return &ShortenedUrlController{
		shortenUrlUseCase:  shortenUrlUseCase,
		redirectUrlUseCase: redirectUrlUseCase,
	}
}

func (uc *ShortenedUrlController) ShortenUrl(c *gin.Context) {
	var req *ports.ShortenUrlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := uc.shortenUrlUseCase.Execute(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, res)
}

func (uc *ShortenedUrlController) RedirectUrl(c *gin.Context) {
	var req *ports.RedirectUrlRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := uc.redirectUrlUseCase.Execute(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(302, res.OriginalUrl)
}
