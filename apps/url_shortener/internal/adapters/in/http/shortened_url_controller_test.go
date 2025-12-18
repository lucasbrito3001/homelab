package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucasbrito3001/url_shortner/internal/app/ports"
	"github.com/lucasbrito3001/url_shortner/internal/app/ports/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUrlController_ShortenUrl(t *testing.T) {
	// mockError := errors.New("some error")
	// mockBody := `{"original_url": "https://google.com"}`
	mockShortenUrlUseCase := mocks.NewShortenUrlUseCase(t)
	mockRedirectUrlUseCase := mocks.NewRedirectUrlUseCase(t)
	controller := NewShortenedUrlController(mockShortenUrlUseCase, mockRedirectUrlUseCase)

	t.Run("should return error 400 with error ErrInvalidRequest when the request is invalid", func(t *testing.T) {
		// Arrange
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		_, engine := gin.CreateTestContext(w)
		engine.POST("/shorten", controller.ShortenUrl)
		invalidBody := `{"invalid_field": 123`
		req, _ := http.NewRequest(http.MethodPost, "/shorten", strings.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")

		// Act
		engine.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error")
		mockShortenUrlUseCase.AssertNotCalled(t, "Execute", mock.Anything, mock.Anything)
	})

	t.Run("should return code successfully", func(t *testing.T) {
		// Arrange
		expectedResp := ports.ShortenUrlResponse{Code: "abc123"}
		mockShortenUrlUseCase.On("Execute", mock.Anything, mock.Anything).Return(&expectedResp, nil)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		body := `{"original_url": "https://google.com"}`
		engine.POST("/shorten", controller.ShortenUrl) // Registra a rota

		req, _ := http.NewRequest(http.MethodPost, "/shorten", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req

		// Act
		engine.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "abc123")
		var actualResponse ports.ShortenUrlResponse
		json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.Equal(t, "abc123", actualResponse.Code)
	})
}
