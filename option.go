package requests

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type ClientOption func(client *Client)

func WithTimeout(timeout time.Duration) ClientOption {
	return func(client *Client) { client.Timeout = timeout }
}

func WithTransport(transport http.RoundTripper) ClientOption {
	return func(client *Client) { client.Transport = transport }
}

func WithCheckRedirect(checkRedirect func(req *http.Request, via []*http.Request) error) ClientOption {
	return func(client *Client) { client.CheckRedirect = checkRedirect }
}

func WithJar(jar http.CookieJar) ClientOption {
	return func(client *Client) { client.Jar = jar }
}

type ReqOption interface {
	Do(req *Request) error
}

type Header map[string]string

func (h Header) Do(req *Request) error {
	for key, value := range h {
		req.Header.Set(key, value)
	}
	return nil
}

type Params map[string]string

func (p Params) Do(req *Request) error {
	if len(p) == 0 {
		return nil
	}
	if req.URL.RawQuery != "" {
		req.URL.RawQuery += "&"
	}
	values := url.Values{}
	for key, value := range p {
		values.Set(key, value)
	}
	req.URL.RawQuery += values.Encode()
	return nil
}

type Json map[string]interface{}

func (j Json) Do(req *Request) error {
	if req.files != nil || req.form != nil {
		return ErrInvalidBodyType
	}
	if req.json == nil {
		req.json = make(Json)
	}
	for k, v := range j {
		req.json[k] = v
	}
	return nil
}

type Jsons []Json

func (j Jsons) Do(req *Request) error {
	if req.files != nil || req.form != nil {
		return ErrInvalidBodyType
	}
	req.jsons = append(req.jsons, j...)
	return nil
}

type Form map[string]string

func (f Form) Do(req *Request) error {
	if req.json != nil || req.jsons != nil {
		return ErrInvalidBodyType
	}
	if req.form == nil {
		req.form = make(Form)
	}
	for k, v := range f {
		req.form[k] = v
	}
	return nil
}

type Cookies map[string]string

func (c Cookies) Do(req *Request) error {
	for name, value := range c {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}
	return nil
}

type file struct {
	field    string // file field
	content  io.ReadCloser
	name     string // file name
	filePath string
}

func FileWithPath(field string, path string) *file {
	return &file{filePath: path, field: field, name: filepath.Base(path)}
}

func FileWithContent(field, fileName string, content []byte) *file {
	return &file{field: field, name: fileName, content: ioutil.NopCloser(bytes.NewReader(content))}
}

func (f file) Do(req *Request) error {
	if f.field == "" {
		return errors.Wrap(ErrInvalidFile, "field is nil")
	} else if f.name == "" {
		return errors.Wrap(ErrInvalidFile, "fileName is nil")
	}
	if f.content == nil && f.filePath == "" {
		return errors.Wrap(ErrInvalidFile, "content is nil, path is nil")
	} else if f.content == nil {
		fc, err := os.Open(f.filePath)
		if err != nil {
			return errors.Wrap(ErrInvalidFile, fmt.Sprintf("open file err: %v", err))
		}
		f.content = fc
	}
	req.files = append(req.files, &f)
	return nil
}
