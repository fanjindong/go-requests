package requests

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ajg/form"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

type Request struct {
	*http.Request
	files []*file
	form  Form
	json  Json
	jsons Jsons
}

// NewRequest wraps NewRequestWithContext using the background context.
func NewRequest(method, url string) (*Request, error) {
	r, err := http.NewRequestWithContext(context.Background(), method, url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", userAgent)
	return &Request{Request: r}, nil
}

func (req *Request) loadBody() error {
	if req.files == nil && req.form == nil && req.json == nil && req.jsons == nil {
		return nil
	}
	// application/json
	var jsonData interface{}
	if req.jsons != nil {
		jsonData = req.jsons
	} else if req.json != nil {
		jsonData = req.json
	}
	if jsonData != nil {
		req.Header.Set("content-Type", "application/json")
		jsonBytes, err := marshal(jsonData)
		if err != nil {
			return errors.Wrap(ErrInvalidJson, err.Error())
		}
		jsonBuffer := bytes.NewBuffer(jsonBytes)
		req.Body = ioutil.NopCloser(jsonBuffer)
		return nil
	}
	// application/x-www-form-urlencoded
	if req.files == nil {
		req.Header.Set("content-Type", "application/x-www-form-urlencoded")
		data, err := form.EncodeToString(req.form)
		if err != nil {
			return errors.Wrap(ErrInvalidForm, err.Error())
		}
		dataReader := strings.NewReader(data)
		req.Body = ioutil.NopCloser(dataReader)
		return nil
	}
	// multipart/form-data; boundary=b...
	buffer := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(buffer)
	for _, file := range req.files {
		writer, err := multipartWriter.CreateFormFile(file.field, file.name)
		if err != nil {
			return errors.Wrap(ErrInvalidFile, fmt.Sprintf("field: %s, name: %s, CreateFormFile err: %v", file.field, file.name, err))
		}
		if _, err = io.Copy(writer, file.content); err != nil && err != io.EOF {
			return errors.Wrap(ErrInvalidFile, fmt.Sprintf("field: %s, name: %s, Copy err: %v", file.field, file.name, err))
		}
		if err = file.content.Close(); err != nil {
			return errors.Wrap(ErrInvalidFile, fmt.Sprintf("field: %s, name: %s, Close err: %v", file.field, file.name, err))
		}
	}
	for k, v := range req.form {
		if err := multipartWriter.WriteField(k, v); err != nil {
			return errors.Wrap(ErrInvalidForm, fmt.Sprintf("Key: %s, Value: %s, WriteField err: %v", k, v, err))
		}
	}
	if err := multipartWriter.Close(); err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(buffer)
	req.Header.Add("content-Type", multipartWriter.FormDataContentType())
	return nil
}
