package requests

import (
	"reflect"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	url := baseUrl + "/timeout"
	tests := []struct {
		input     []Option
		wantError bool
	}{
		{input: []Option{Timeout(4 * time.Second)}, wantError: false},
		{input: []Option{Timeout(3100 * time.Millisecond)}, wantError: false},
		{input: []Option{Timeout(3000 * time.Millisecond)}, wantError: true},
	}

	for _, tt := range tests {
		_, err := Get(url, tt.input...)
		if !reflect.DeepEqual(err != nil, tt.wantError) {
			t.Errorf("Get() err = %v, wantError %v", err, tt.wantError)
		}
	}
}

//func TestCookies(t *testing.T) {
//	url := baseUrl
//	tests := []struct {
//		input []Option
//		want  map[string]interface{}
//	}{
//		{input: []Option{Cookies{"name": "fjd"}}, want: map[string]interface{}{"name": "fjd"}},
//		{input: []Option{Cookies{"name": "fjd"}, Cookies{"age": "18"}}, want: map[string]interface{}{"name": "fjd", "age": "18"}},
//		{input: []Option{Json{"a": 2.1}, Cookies{"name": "fjd"}, Cookies{"age": "18"}}, want: map[string]interface{}{"a": 2.1, "name": "fjd", "age": "18"}},
//	}
//
//	for _, tt := range tests {
//		resp, err := Post(url, tt.input...)
//		assert.NoError(t, err)
//		respData := make(map[string]interface{})
//		err = resp.Json(&respData)
//		assert.NoError(t, err)
//		got := respData["data"]
//		assert.EqualValues(t, tt.want, got)
//	}
//}

func TestOptions(t *testing.T) {
	type args struct {
		method  string
		url     string
		options []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    map[string]interface{}
	}{
		{name: "1", args: args{method: "GET", url: baseUrl + "/get", options: []Option{Params{"a": "1"}}}, wantErr: false, want: map[string]interface{}{"a": "1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := session.Request(tt.args.method, tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := map[string]interface{}{}
			if err = response.Json(&got); (err != nil) != tt.wantErr {
				t.Errorf("response.Json() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
