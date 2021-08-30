package requests

import (
	"reflect"
	"testing"
	"time"
)

func TestWithTimeout(t *testing.T) {
	url := testUrl + "/timeout"
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
	}{
		{args: args{}},
		{args: args{timeout: 2 * time.Second}},
		{args: args{timeout: 1100 * time.Millisecond}},
		{args: args{timeout: 1000 * time.Millisecond}, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(WithTimeout(tt.args.timeout))
			_, err := client.Get(url)
			if !reflect.DeepEqual(err != nil, tt.wantError) {
				t.Errorf("WithTimeout() err = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestHeaders(t *testing.T) {
	url := testUrl + "/header"
	type args struct {
		headers []ReqOption
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{args: args{}, want: map[string][]string{"Accept-Encoding": {"gzip"}, "User-Agent": {userAgent}}},
		{args: args{headers: []ReqOption{Header{"a": "1"}}}, want: map[string][]string{"Accept-Encoding": {"gzip"}, "User-Agent": {userAgent}, "A": {"1"}}},
		{args: args{headers: []ReqOption{Header{"a": "1", "b": "2"}}}, want: map[string][]string{"Accept-Encoding": {"gzip"}, "User-Agent": {userAgent}, "A": {"1"}, "B": {"2"}}},
		{args: args{headers: []ReqOption{Header{"a": "1"}, Header{"b": "2"}}}, want: map[string][]string{"Accept-Encoding": {"gzip"}, "User-Agent": {userAgent}, "A": {"1"}, "B": {"2"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := Get(url, tt.args.headers...)
			got := make(map[string][]string)
			_ = resp.Json(&got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Header() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams(t *testing.T) {
	url := testUrl + "/get"
	type args struct {
		opts []ReqOption
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{args: args{}, want: map[string]string{}},
		{args: args{opts: []ReqOption{Params{"a": "1"}}}, want: map[string]string{"a": "1"}},
		{args: args{opts: []ReqOption{Params{"a": "1", "b": "2"}}}, want: map[string]string{"a": "1", "b": "2"}},
		{args: args{opts: []ReqOption{Params{"a": "1"}, Params{"b": "2"}}}, want: map[string]string{"a": "1", "b": "2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := Get(url, tt.args.opts...)
			got := make(map[string]string)
			_ = resp.Json(&got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Params() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson(t *testing.T) {
	url := testUrl + "/post"
	type args struct {
		opts []ReqOption
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{args: args{}, want: map[string]interface{}{}},
		{args: args{opts: []ReqOption{Json{"a": "1"}}}, want: map[string]interface{}{"a": "1"}},
		{args: args{opts: []ReqOption{Json{"a": "1", "b": 2}}}, want: map[string]interface{}{"a": "1", "b": float64(2)}},
		{args: args{opts: []ReqOption{Json{"a": "1"}, Json{"b": "2"}}}, want: map[string]interface{}{"b": "2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := Post(url, tt.args.opts...)
			got := make(map[string]interface{})
			_ = resp.Json(&got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Json() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsons(t *testing.T) {
	url := testUrl + "/post"
	type args struct {
		opts []ReqOption
	}
	tests := []struct {
		name string
		args args
		want []map[string]interface{}
	}{
		{args: args{}, want: []map[string]interface{}{}},
		{args: args{opts: []ReqOption{Jsons{{"a": "1"}}}}, want: []map[string]interface{}{{"a": "1"}}},
		{args: args{opts: []ReqOption{Jsons{{"a": "1"}, {"b": 2}}}}, want: []map[string]interface{}{{"a": "1"}, {"b": float64(2)}}},
		{args: args{opts: []ReqOption{Jsons{{"a": "1", "b": 2}}}}, want: []map[string]interface{}{{"a": "1", "b": float64(2)}}},
		{args: args{opts: []ReqOption{Jsons{{"a": "1"}}, Jsons{{"b": "2"}}}}, want: []map[string]interface{}{{"b": "2"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := Post(url, tt.args.opts...)
			got := make([]map[string]interface{}, 0)
			_ = resp.Json(&got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Jsons() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForm(t *testing.T) {
	url := testUrl + "/post"
	type args struct {
		opts []ReqOption
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{args: args{}, want: map[string]interface{}{}},
		{args: args{opts: []ReqOption{Form{"a": "1"}}}, want: map[string]interface{}{"a": "1"}},
		{args: args{opts: []ReqOption{Form{"a": "1", "b": "2"}}}, want: map[string]interface{}{"a": "1", "b": "2"}},
		{args: args{opts: []ReqOption{Form{"a": "1"}, Form{"b": "2"}}}, want: map[string]interface{}{"b": "2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := Post(url, tt.args.opts...)
			got := make(map[string]interface{})
			_ = resp.Json(&got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Form() got = %v, want %v", got["b"].(int), tt.want)
			}
		})
	}
}
