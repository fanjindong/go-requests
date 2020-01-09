package requests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var BaseUrl = "http://127.0.0.1:8080"

func TestGet(t *testing.T) {
	type input struct {
		url    string
		params Params
	}

	tests := []struct {
		input input
		want  string
	}{
		{input: input{url: BaseUrl + "/get", params: Params{"a": "1"}}, want: "/get?a=1"},
		{input: input{url: BaseUrl + "/get", params: Params{"a": "1", "b": "2"}}, want: "/get?a=1&b=2"},
		{input: input{url: BaseUrl + "/get?", params: Params{"a": "1"}}, want: "/get?a=1"},
		{input: input{url: BaseUrl + "/get?a=1", params: Params{}}, want: "/get?a=1"},
		{input: input{url: BaseUrl + "/get?a=1", params: Params{"b": "2"}}, want: "/get?a=1&b=2"},
	}

	for _, ts := range tests {
		resp, err := Get(ts.input.url, ts.input.params)
		assert.NoError(t, err)
		respData := make(map[string]interface{})
		err = resp.Json(&respData)
		assert.NoError(t, err)
		got := respData["url"]
		assert.Equal(t, ts.want, got)
	}
}

func TestPost(t *testing.T) {
	url := BaseUrl + "/post"
	tests := []struct {
		input []Option
		want  map[string]interface{}
	}{
		{input: []Option{Json{"a": 1.1}}, want: map[string]interface{}{"a": 1.1}},
		{input: []Option{Params{"params": "1"}, Json{"b": 2.2}}, want: map[string]interface{}{"params": "1", "b": 2.2}},
		{input: []Option{Params{"params": "1"}, Json{"b": 2.2}, Headers{"headers": "22"}}, want: map[string]interface{}{"params": "1", "b": 2.2, "headers": "22"}},
	}

	for _, ts := range tests {
		resp, err := Post(url, ts.input...)
		assert.NoError(t, err)
		respData := make(map[string]interface{})
		err = resp.Json(&respData)
		assert.NoError(t, err)
		got := respData["data"]
		assert.EqualValues(t, ts.want, got)
	}
}

func TestPut(t *testing.T) {
	url := BaseUrl + "/put"
	tests := []struct {
		input []Option
		want  map[string]interface{}
	}{
		{input: []Option{Json{"a": 1.1}}, want: map[string]interface{}{"a": 1.1}},
		{input: []Option{Params{"params": "1"}, Json{"b": 2.2}}, want: map[string]interface{}{"params": "1", "b": 2.2}},
		{input: []Option{Params{"params": "1"}, Data{"form": "2.2"}}, want: map[string]interface{}{"params": "1", "form": "2.2"}},
	}

	for _, ts := range tests {
		resp, err := Put(url, ts.input...)
		assert.NoError(t, err)
		respData := make(map[string]interface{})
		err = resp.Json(&respData)
		assert.NoError(t, err)
		got := respData["data"]
		assert.EqualValues(t, ts.want, got)
	}
}

func TestTimeout(t *testing.T) {
	url := BaseUrl + "/timeout"
	tests := []struct {
		input []Option
		want  int
	}{
		{input: []Option{Timeout(4 * time.Second)}, want: 200},
		{input: []Option{Timeout(3100 * time.Millisecond)}, want: 200},
		{input: []Option{Timeout(3000 * time.Millisecond)}, want: 503},
	}

	for _, ts := range tests {
		resp, err := Get(url, ts.input...)
		if ts.want != 200 {
			assert.True(t, err != nil)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, ts.want, resp.StatusCode)
		}
	}
}

func TestCookies(t *testing.T) {
	url := BaseUrl + "/post"
	tests := []struct {
		input []Option
		want  map[string]interface{}
	}{
		{input: []Option{Cookies{Name: "name", Value: "fjd"}}, want: map[string]interface{}{"name": "fjd"}},
		{input: []Option{Cookies{Name: "name", Value: "fjd"}, Cookies{Name: "age", Value: "18"}}, want: map[string]interface{}{"name": "fjd", "age": "18"}},
		{input: []Option{Json{"a": 2.1}, Cookies{Name: "name", Value: "fjd"}, Cookies{Name: "age", Value: "18"}}, want: map[string]interface{}{"a": 2.1, "name": "fjd", "age": "18"}},
	}

	for _, ts := range tests {
		resp, err := Post(url, ts.input...)
		assert.NoError(t, err)
		respData := make(map[string]interface{})
		err = resp.Json(&respData)
		assert.NoError(t, err)
		got := respData["data"]
		assert.EqualValues(t, ts.want, got)
	}
}

func TestFiles(t *testing.T) {
	_, err := FileFromPath("./go.mod")
	assert.NoError(t, err)
	_ = FileFromContents("demo.text", "123 \n")
}
