package requests

import "github.com/pkg/errors"

var (
	ErrInvalidJson = errors.New("invalid Json value")

	// ErrUnrecognizedEncoding will be throwed while changing response encoding
	// if encoding is not recognized
	ErrUnrecognizedEncoding = errors.New("Unrecognized encoding")

	// ErrInvalidMethod will be throwed when method not in
	// [HEAD, GET, POST, DELETE, OPTIONS, PUT, PATCH, CONNECT, TRACE]
	ErrInvalidMethod = errors.New("Method is invalid")
)
