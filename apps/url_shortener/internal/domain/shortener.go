package domain

import (
	"github.com/jxskiss/base62"
)

type Shortener struct {
	encoding *base62.Encoding
}

func NewShortener(alphabet string) (*Shortener, error) {
	if len(alphabet) != 62 {
		return nil, ErrInvalidAlphabet
	}

	return &Shortener{
		encoding: base62.NewEncoding(alphabet),
	}, nil
}

func (s *Shortener) Encode(id int64) string {
	return string(s.encoding.FormatInt(id))
}
