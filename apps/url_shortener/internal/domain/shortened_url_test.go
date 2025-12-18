package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShortenedUrl(t *testing.T) {
	t.Run("should return a new url with a ID and CreatedAt", func(t *testing.T) {
		url := NewShortenedUrl("original", "code")

		assert.NotEmpty(t, url.ID)
		assert.NotEmpty(t, url.CreatedAt)
		assert.Equal(t, len(url.ID), 36)
		assert.Equal(t, url.Code, "code")
		assert.Equal(t, url.Original, OriginalUrl("original"))
	})
}

func TestNewOriginalUrl(t *testing.T) {
	t.Run("should return ErrNotSecureOriginalUrl when the url is not secure", func(t *testing.T) {
		// Act
		_, err := NewOriginalUrl("http://www.mock.com")

		// Assert
		assert.ErrorIs(t, err, ErrNotSecureOriginalUrl)
	})

	t.Run("should return error when the url is invalid", func(t *testing.T) {
		// Act
		_, err := NewOriginalUrl("invalid")

		// Assert
		assert.NotNil(t, err)
	})

	t.Run("should return a OriginalUrl when the url is valid", func(t *testing.T) {
		// Act
		url, err := NewOriginalUrl("https://www.mock.com")

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, url, OriginalUrl("https://www.mock.com"))
	})
}
