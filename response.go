package requests

import (
	"io/ioutil"
	"net/http"
	"os"
)

// Response is the wrapper for http.Response
type Response struct {
	*http.Response
	bytes []byte
}

func NewResponse(r *http.Response) (*Response, error) {
	resp := &Response{Response: r}
	_, err := resp.Bytes()
	_ = r.Body.Close()
	return resp, err
}

func (r *Response) Text() string {
	return string(r.bytes)
}

func (r *Response) Bytes() ([]byte, error) {
	if r.bytes == nil {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		r.bytes = data
	}
	return r.bytes, nil
}

// Json could parse http json response
func (r Response) Json(s interface{}) error {
	return unmarshal(r.bytes, s)
}

// SaveFile save bytes data to a local file
func (r Response) SaveFile(filename string) error {
	dst, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = dst.Close()
	}()

	if bt, err := r.Bytes(); err != nil {
		return err
	} else {
		_, err = dst.Write(bt)
		return err
	}
}
