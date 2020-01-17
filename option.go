package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ajg/form"
	"github.com/pkg/errors"
)

type Option interface {
	ApplyClient(client *http.Client)
	ApplyRequest(req *http.Request) error
}

type Headers map[string]string

func (h Headers) ApplyClient(_ *http.Client) {}
func (h Headers) ApplyRequest(req *http.Request) error {
	for key, value := range h {
		req.Header.Set(key, value)
	}
	return nil
}

type Timeout time.Duration

func (t Timeout) ApplyClient(client *http.Client) {
	client.Timeout = time.Duration(t)
}
func (t Timeout) ApplyRequest(_ *http.Request) error { return nil }

type Params map[string]interface{}

func (p Params) ApplyClient(_ *http.Client) {}
func (p Params) ApplyRequest(req *http.Request) error {
	var rawQuery []string
	if req.URL.RawQuery != "" {
		rawQuery = append(rawQuery, req.URL.RawQuery)
	}

	for key, value := range p {
		rawQuery = append(rawQuery, fmt.Sprintf("%s=%s", key, value))
	}
	req.URL.RawQuery = strings.Join(rawQuery, "&")
	return nil
}

type Json map[string]interface{}

func (j Json) ApplyClient(_ *http.Client) {}
func (j Json) ApplyRequest(req *http.Request) error {
	jsonBytes, err := json.Marshal(j)
	if err != nil {
		return errors.Wrap(ErrInvalidJson, err.Error())
	}
	jsonBuffer := bytes.NewBuffer(jsonBytes)
	req.Body = ioutil.NopCloser(jsonBuffer)
	req.Header.Set("Content-Type", "application/json")
	return nil
}

type Data map[string]interface{}

func (d Data) ApplyClient(_ *http.Client) {}
func (d Data) ApplyRequest(req *http.Request) error {
	data, err := form.EncodeToString(d)
	if err != nil {
		return errors.Wrap(ErrInvalidForm, err.Error())
	}
	dataReader := strings.NewReader(data)
	req.Body = ioutil.NopCloser(dataReader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return nil
}

// File is a struct that is used to specify the file that a User wishes to upload.
type File struct {
	// Filename is the name of the file that you wish to upload. We use this to guess the mimetype as well as pass it onto the server
	FileName string
	// FileContent is happy as long as you pass it a io.ReadCloser (which most file use anyways)
	FileContent io.ReadCloser
	// MimeType represents which mimetime should be sent along with the file.
	// When empty, defaults to application/octet-stream
	MimeType string
}

// FName changes file's filename in multipart form
// invoke it in a chain
func (f *File) FName(filename string) *File {
	f.FileName = filename
	return f
}

// MIME changes file's mime type in multipart form
// invoke it in a chain
func (f *File) MIME(mimeType string) *File {
	f.MimeType = mimeType
	return f
}

// File returns a new file struct
func FileFromContents(filename string, content string) *File {
	return &File{FileContent: ioutil.NopCloser(strings.NewReader(content)), FileName: filename}
}

// FileFromPath returns a file struct from file path
func FileFromPath(filePath string) (*File, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidFile, err.Error())
	}
	return &File{FileContent: fd, FileName: filepath.Base(filePath)}, nil
}

type Files map[string]interface{}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string { return quoteEscaper.Replace(s) }

func (f Files) ApplyClient(_ *http.Client) {}
func (f Files) ApplyRequest(req *http.Request) error {
	buffer := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(buffer)

	for key, value := range f {
		switch value := value.(type) {
		case *File:
			if value.FileContent == nil || value.FileName == "" {
				return ErrInvalidFile
			}
			var writer io.Writer
			var err error

			if value.MimeType == "" {
				writer, err = multipartWriter.CreateFormFile(key, value.FileName)
			} else {
				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes(key), escapeQuotes(value.FileName)))
				h.Set("Content-Type", value.MimeType)
				writer, err = multipartWriter.CreatePart(h)
			}
			if err != nil {
				return errors.Wrap(ErrInvalidFile, err.Error())
			}

			if _, err = io.Copy(writer, value.FileContent); err != nil && err != io.EOF {
				return errors.Wrap(ErrInvalidFile, err.Error())
			}

			if err := value.FileContent.Close(); err != nil {
				return errors.Wrap(ErrInvalidFile, err.Error())
			}
		case string:
			err := multipartWriter.WriteField(key, value)
			if err != nil {
				return errors.Wrap(ErrInvalidFile, err.Error())
			}
		default:
			return errors.Wrap(ErrInvalidFile, fmt.Sprintf("invalid value: %+v", value))
		}
	}

	if err := multipartWriter.Close(); err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(buffer)
	req.Header.Add("Content-Type", multipartWriter.FormDataContentType())
	return nil
}

type Cookies map[string]string

func (c Cookies) ApplyClient(_ *http.Client) {}
func (c Cookies) ApplyRequest(req *http.Request) error {
	for name, value := range c {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}
	return nil
}
