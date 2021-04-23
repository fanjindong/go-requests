package requests

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var (
	baseUrl = "http://127.0.0.1:8080"
	session *Session
)

func TestGet(t *testing.T) {
	type input struct {
		url    string
		params Params
	}

	tests := []struct {
		input input
		want  map[string]string
	}{
		{input: input{url: baseUrl + "/get", params: Params{"a": "1"}}, want: map[string]string{"a": "1"}},
		{input: input{url: baseUrl + "/get", params: Params{"a": "1", "b": "2"}}, want: map[string]string{"a": "1", "b": "2"}},
		{input: input{url: baseUrl + "/get?"}, want: map[string]string{}},
		{input: input{url: baseUrl + "/get?a=1", params: Params{}}, want: map[string]string{"a": "1"}},
		{input: input{url: baseUrl + "/get?a=1", params: Params{"b": "2"}}, want: map[string]string{"a": "1", "b": "2"}},
	}

	for _, tt := range tests {
		r, err := Get(tt.input.url, tt.input.params)
		if err != nil {
			panic(err)
		}
		got := map[string]string{}
		if err = r.Json(&got); err != nil {
			t.Log(r.Text())
			panic(err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Get() got = %v, want %v", got, tt.want)
		}
	}
}

func TestPost(t *testing.T) {
	url := baseUrl + "/post"
	tests := []struct {
		input []Option
		want  map[string]interface{}
	}{
		{input: []Option{Json{"a": 1.1}}, want: map[string]interface{}{"a": 1.1}},
		{input: []Option{Json{"a": 1.1, "b": 2.2}}, want: map[string]interface{}{"a": 1.1, "b": 2.2}},
		{input: []Option{Data{"a": 1.1, "b": 2.2}}, want: map[string]interface{}{"a": "1.1", "b": "2.2"}},
	}

	for _, tt := range tests {
		resp, err := Post(url, tt.input...)
		if err != nil {
			panic(err)
		}
		got := make(map[string]interface{})
		if err = resp.Json(&got); err != nil {
			panic(err)
		}
		assert.EqualValues(t, tt.want, got)
	}
}

func TestFiles(t *testing.T) {
	_, err := FilePath("./go.mod")
	if err != nil {
		panic(err)
	}
	_ = FileContents("demo.text", "123 \n")
}

func BenchmarkGetRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Get(baseUrl)
		if err != nil {
			panic(err)
		}
		//resp.Text()
	}
}
