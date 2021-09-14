package requests

import (
	"encoding/json"
	"net/http"
	"strings"
)

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

func Get(url string, opts ...ReqOption) (*Response, error) {
	return DefaultClient.Request(GET, url, opts...)
}

func Post(url string, opts ...ReqOption) (*Response, error) {
	return DefaultClient.Request(POST, url, opts...)
}

func Put(url string, opts ...ReqOption) (*Response, error) {
	return DefaultClient.Request(PUT, url, opts...)
}

func Delete(url string, opts ...ReqOption) (*Response, error) {
	return DefaultClient.Request(DELETE, url, opts...)
}

func ReqOptions(url string, opts ...ReqOption) (*Response, error) {
	return DefaultClient.Request(OPTIONS, url, opts...)
}

func Patch(url string, opts ...ReqOption) (*Response, error) {
	return DefaultClient.Request(PATCH, url, opts...)
}

func Head(url string, opts ...ReqOption) (*Response, error) {
	return DefaultClient.Request(HEAD, url, opts...)
}

type Client struct {
	*http.Client
	hooks []Hook
}

var DefaultClient = &Client{Client: http.DefaultClient}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{Client: &http.Client{}}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (s *Client) Request(method, url string, opts ...ReqOption) (*Response, error) {
	method = strings.ToUpper(method)
	switch method {
	case HEAD, GET, POST, DELETE, OPTIONS, PUT, PATCH:
	default:
		return nil, ErrInvalidMethod
	}

	req, err := NewRequest(method, url)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		err = opt.Do(req)
		if err != nil {
			return nil, err
		}
	}
	if err = req.loadBody(); err != nil {
		return nil, err
	}

	for _, h := range s.hooks {
		h.BeforeProcess(req)
	}
	var result *http.Response
	var resp *Response
	success := make(chan struct{})
	done := req.Context().Done()
	if done != nil {
		go func() {
			result, err = s.Do(req.Request)
			close(success)
		}()
		select {
		case <-done:
			err = ErrTimeout
		case <-success:
		}
	} else {
		result, err = s.Do(req.Request)
	}
	if err == nil {
		resp, err = NewResponse(result)
	}
	for _, h := range s.hooks {
		h.AfterProcess(req, resp, err)
	}
	return resp, err
}

func (s *Client) Get(url string, opts ...ReqOption) (*Response, error) {
	return s.Request(GET, url, opts...)
}

func (s *Client) Post(url string, opts ...ReqOption) (*Response, error) {
	return s.Request(POST, url, opts...)
}

func (s *Client) Put(url string, opts ...ReqOption) (*Response, error) {
	return s.Request(PUT, url, opts...)
}

func (s *Client) Delete(url string, opts ...ReqOption) (*Response, error) {
	return s.Request(DELETE, url, opts...)
}

func (s *Client) ReqOptions(url string, opts ...ReqOption) (*Response, error) {
	return s.Request(OPTIONS, url, opts...)
}

func (s *Client) Patch(url string, opts ...ReqOption) (*Response, error) {
	return s.Request(PATCH, url, opts...)
}

func (s *Client) Head(url string, opts ...ReqOption) (*Response, error) {
	return s.Request(HEAD, url, opts...)
}

func (s *Client) AddHook(h Hook) {
	s.hooks = append(s.hooks, h)
}

var unmarshal = json.Unmarshal
var marshal = json.Marshal

////SetUnmarshal Set custom Unmarshal functions, default is json.Unmarshal
//func SetUnmarshal(f func(data []byte, v interface{}) error) {
//	unmarshal = f
//}
//
////SetMarshal Set custom Marshal functions, default is json.Marshal
//func SetMarshal(f func(v interface{}) ([]byte, error)) {
//	marshal = f
//}
