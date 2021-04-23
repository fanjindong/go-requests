package requests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/axgle/mahonia"
)

// Response is the wrapper for http.Response
type Response struct {
	*http.Response
	encoding string
	bytes    []byte
	Headers  *http.Header
}

func NewResponse(r *http.Response) (*Response, error) {
	resp := &Response{
		Response: r,
		encoding: "utf-8",
		Headers:  &r.Header,
	}
	_, err := resp.Bytes()
	_ = r.Body.Close()
	return resp, err
}

func (r *Response) Text() (string, error) {
	if bt, err := r.Bytes(); err != nil {
		return "", err
	} else {
		return string(bt), nil
	}
}

func (r *Response) Bytes() ([]byte, error) {
	if r.bytes == nil {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		// for multiple reading
		// e.g. goquery.NewDocumentFromReader
		//r.Body = ioutil.NopCloser(bytes.NewBuffer(data))

		if r.encoding != "utf-8" {
			data = []byte(mahonia.NewDecoder(r.encoding).ConvertString(string(data)))
		}
		r.bytes = data
	}

	return r.bytes, nil
}

// Json could parse http json response
func (r Response) Json(s interface{}) error {
	// Json response not must be `application/json` type
	// maybe `text/plain`...etc.
	// requests will parse it regardless of the content-type
	/*
		cType := r.Header.Get("Content-Type")
		if !strings.Contains(cType, "json") {
			return ErrNotJsonResponse
		}
	*/
	if bt, err := r.Bytes(); err != nil {
		return err
	} else {
		return json.Unmarshal(bt, s)
	}
}

// SetEncode changes Response.encoding
// and it changes Response.Text every times be invoked
func (r *Response) SetEncode(e string) error {
	if r.encoding != e {
		if mahonia.NewDecoder(e) == nil {
			return ErrUnrecognizedEncoding
		}
		r.encoding = strings.ToLower(e)
		if r.bytes != nil {
			r.bytes = []byte(mahonia.NewDecoder(r.encoding).ConvertString(string(r.bytes)))
		}
	}
	return nil
}

// GetEncode returns Response.encoding
func (r Response) GetEncode() string {
	return r.encoding
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
