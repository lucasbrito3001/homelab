package domain

import "errors"

var (
	ErrInvalidAlphabet      = errors.New("the alphabet must have 62 characters")
	ErrEmptyOriginalUrl     = errors.New("the original url cannot be empty")
	ErrNotSecureOriginalUrl = errors.New("the original url must be secure (https)")
	ErrInvalidOriginalUrl   = errors.New("the original url is invalid, must be in format https://{subdomain}.{domain}.{top-level-domain}")
	ErrTooLargeCode         = errors.New("the code cannot be longer than 7 characters")
	ErrNonAlphanumericCode  = errors.New("the code must be alphanumeric")
)
