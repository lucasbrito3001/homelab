package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortener_Encode(t *testing.T) {
	shortener, _ := NewShortener("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	t.Run("should return 24 when the id is 128", func(t *testing.T) {
		code := shortener.Encode(128)

		assert.Equal(t, code, "24")
	})

}

func TestShortener_NewShortener(t *testing.T) {
	invalidAlphabet := "123"

	t.Run("should return error when the alphabet dont have 62 characters", func(t *testing.T) {
		_, err := NewShortener(invalidAlphabet)

		assert.ErrorIs(t, err, ErrInvalidAlphabet)
	})
}
