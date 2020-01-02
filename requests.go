package requests

const (
	version   = "0.0.1"
	userAgent = "go-requests/" + version
	author    = "fanjindong"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	HEAD    = "HEAD"
)

func Get(url string, option ...Option) (*Response, error) {
	s := NewSession()
	return s.Request(GET, url, option...)
}

func Post(url string, option ...Option) (*Response, error) {
	s := NewSession()
	return s.Request(POST, url, option...)
}

func Put(url string, option ...Option) (*Response, error) {
	s := NewSession()
	return s.Request(PUT, url, option...)
}

func Delete(url string, option ...Option) (*Response, error) {
	s := NewSession()
	return s.Request(DELETE, url, option...)
}

func Options(url string, option ...Option) (*Response, error) {
	s := NewSession()
	return s.Request(OPTIONS, url, option...)
}

func Patch(url string, option ...Option) (*Response, error) {
	s := NewSession()
	return s.Request(PATCH, url, option...)
}

func Head(url string, option ...Option) (*Response, error) {
	s := NewSession()
	return s.Request(HEAD, url, option...)
}
