package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/lucasbrito3001/url_shortner/internal/app/ports"
	"github.com/lucasbrito3001/url_shortner/internal/app/ports/mocks"
	"github.com/lucasbrito3001/url_shortner/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenUrl_Execute(t *testing.T) {
	mockOriginalUrl := "https://www.mock.com"
	someError := errors.New("some error")
	mockUrlRepository := mocks.NewUrlRepository(t)
	mockCounterRepository := mocks.NewCounterRepository(t)
	shortener, _ := domain.NewShortener("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	usecase := NewShortenUrl(mockUrlRepository, mockCounterRepository, shortener, 0)

	t.Run("should return error when Increment fails", func(t *testing.T) {
		// Arrange
		mockCounterRepository.On("Increment", mock.Anything).Once().Return(int64(0), someError)

		// Act
		response, err := usecase.Execute(context.Background(), &ports.ShortenUrlRequest{
			OriginalUrl: mockOriginalUrl,
		})

		// Assert
		assert.ErrorIs(t, err, someError)
		assert.Nil(t, response)
	})

	t.Run("should return error when Save fails", func(t *testing.T) {
		// Arrange
		mockCounterRepository.On("Increment", mock.Anything).Once().Return(int64(0), nil)
		mockUrlRepository.On("Save", mock.Anything, mock.Anything).Once().Return(someError)

		// Act
		response, err := usecase.Execute(context.Background(), &ports.ShortenUrlRequest{
			OriginalUrl: mockOriginalUrl,
		})

		// Assert
		assert.ErrorIs(t, err, someError)
		assert.Nil(t, response)
	})

	t.Run("should return a shorten url successfully", func(t *testing.T) {
		// Arrange
		mockUrlRepository.On("Save", mock.Anything, mock.Anything).Once().Return(nil)
		mockCounterRepository.On("Increment", mock.Anything).Once().Return(int64(1), nil)

		// Act
		response, err := usecase.Execute(context.Background(), &ports.ShortenUrlRequest{
			OriginalUrl: mockOriginalUrl,
		})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, response.Code, "1")
		assert.IsType(t, response, &ports.ShortenUrlResponse{})
	})
}
