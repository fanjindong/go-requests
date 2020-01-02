package requests

import (
	"errors"
	"net/http"
	"strings"
	"sync"
)

type Session struct {
	client *http.Client
	req    *http.Request
	sync.Mutex
}

func NewSession() *Session {
	client := &http.Client{}
	return &Session{client: client}
}

func (s *Session) Request(method, url string, option ...Option) (*Response, error) {
	s.Lock()
	defer s.Unlock()

	method = strings.ToUpper(method)
	switch method {
	case HEAD, GET, POST, DELETE, OPTIONS, PUT, PATCH:
	default:
		return nil, ErrInvalidMethod
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Close = true

	for _, opt := range option {
		opt.ApplyClient(s.client)
		err = opt.ApplyRequest(req)
		if err != nil {
			return nil, err
		}
	}
	defer func() { // all client config will be restored to the default value after every request
		s.client.CheckRedirect = defaultCheckRedirect
		s.client.Timeout = 0
		s.client.Transport = &http.Transport{}
	}()

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	return NewResponse(resp)
}

// http's defaultCheckRedirect
func defaultCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return errors.New("stopped after 10 redirects")
	}
	return nil
}

func (s *Session) GetRequest() *http.Request {
	return s.req
}

func (s *Session) Get(url string, option ...Option) (*Response, error) {
	return s.Request(GET, url, option...)
}

func (s *Session) Post(url string, option ...Option) (*Response, error) {
	return s.Request(POST, url, option...)
}

func (s *Session) Put(url string, option ...Option) (*Response, error) {
	return s.Request(PUT, url, option...)
}

func (s *Session) Delete(url string, option ...Option) (*Response, error) {
	return s.Request(DELETE, url, option...)
}

func (s *Session) Options(url string, option ...Option) (*Response, error) {
	return s.Request(OPTIONS, url, option...)
}

func (s *Session) Patch(url string, option ...Option) (*Response, error) {
	return s.Request(PATCH, url, option...)
}

func (s *Session) Head(url string, option ...Option) (*Response, error) {
	return s.Request(HEAD, url, option...)
}
