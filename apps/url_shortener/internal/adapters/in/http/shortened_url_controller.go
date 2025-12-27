package http

import (
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := uc.shortenUrlUseCase.Execute(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (uc *ShortenedUrlController) RedirectUrl(c *gin.Context) {
	var req *ports.RedirectUrlRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := uc.redirectUrlUseCase.Execute(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, res.OriginalUrl)
}
