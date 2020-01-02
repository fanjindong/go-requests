package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ajg/form"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type Option interface {
	ApplyClient(client *http.Client)
	ApplyRequest(req *http.Request) error
}

type Params map[string]interface{}

func (p Params) ApplyClient(client *http.Client) {}

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

func (j Json) ApplyClient(client *http.Client) {}
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

func (d Data) ApplyClient(client *http.Client) {}
func (d Data) ApplyRequest(req *http.Request) error {
	data, err := form.EncodeToString(d)
	if err != nil {
		return errors.Wrap(ErrInvalidJson, err.Error())
	}
	dataReader := strings.NewReader(data)
	req.Body = ioutil.NopCloser(dataReader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return nil
}
