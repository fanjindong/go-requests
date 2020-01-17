package requests

import "github.com/pkg/errors"

var (
	// ErrInvalidJson will be throw out when request form body data can not be Marshal
	ErrInvalidForm = errors.New("go-requests: Invalid Form value")

	// ErrInvalidJson will be throw out when request json body data can not be Marshal
	ErrInvalidJson = errors.New("go-requests: Invalid Json value")

	// ErrUnrecognizedEncoding will be throw out while changing response encoding
	// if encoding is not recognized
	ErrUnrecognizedEncoding = errors.New("go-requests: Unrecognized encoding")

	// ErrInvalidMethod will be throw out when method not in
	// [HEAD, GET, POST, DELETE, OPTIONS, PUT, PATCH, CONNECT, TRACE]
	ErrInvalidMethod = errors.New("go-requests: Method is invalid")

	//ErrInvalidFile will be throw out when get file content data
	ErrInvalidFile = errors.New("go-requests: Invalid File content")
)
