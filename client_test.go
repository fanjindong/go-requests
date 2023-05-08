package requests

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var testUrl = fmt.Sprintf("http://127.0.0.1:%d", port)

func TestGet(t *testing.T) {
	type args struct {
		url    string
		option []ReqOption
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{name: "1", args: args{url: testUrl + "/get", option: []ReqOption{}}, want: map[string]string{}},
		{name: "2", args: args{url: testUrl + "/get", option: []ReqOption{Params{"a": "1"}}}, want: map[string]string{"a": "1"}},
		{name: "3", args: args{url: testUrl + "/get", option: []ReqOption{Params{"a": "1", "b": "x"}}}, want: map[string]string{"a": "1", "b": "x"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := Get(tt.args.url, tt.args.option...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := map[string]string{}
			if err := resp.Json(&got); err != nil {
				t.Errorf("resp.Json() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPost(t *testing.T) {
	type args struct {
		url    string
		option []ReqOption
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{name: "1", args: args{url: testUrl + "/post", option: []ReqOption{Json{}}}, want: map[string]interface{}{}},
		{name: "2", args: args{url: testUrl + "/post", option: []ReqOption{Json{"a": "1"}}}, want: map[string]interface{}{"a": "1"}},
		{name: "3", args: args{url: testUrl + "/post", option: []ReqOption{Json{"a": "x", "b": 1.2}}}, want: map[string]interface{}{"a": "x", "b": 1.2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := Post(tt.args.url, tt.args.option...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := map[string]interface{}{}
			if err := resp.Json(&got); err != nil {
				t.Errorf("resp.Json() error = %v, text = %v", err, resp.Text())
				return
			}
			//for k := range got {
			//	if got[k] != tt.want[k] {
			//		t.Errorf("Post() k = %v, got = %+v, want %+v", k, got[k].(int), tt.want[k].(int))
			//	}
			//}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Post() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReuseConnection(t *testing.T) {
	for i := 0; i < 10; i++ {
		resp, err := Get(testUrl)
		if err != nil {
			t.Log(err)
			return
		}
		//data, _ := ioutil.ReadAll(resp.Body)
		t.Log(resp.Status, resp.Text())
		//resp.Body.Close()
		time.Sleep(1000 * time.Millisecond)
	}
}

//BenchmarkGetRequest-8   	   27895	     43212 ns/op
func BenchmarkGetRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Get(testUrl)
		if err != nil {
			panic(err)
		}
		//resp.Text()
	}
}

//
//func TestSetUnmarshalAndSetMarshal(t *testing.T) {
//	type args struct {
//		v         interface{}
//		unmarshal func(data []byte, v interface{}) error
//		marshal   func(v interface{}) ([]byte, error)
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{args: args{v: map[string]interface{}{"a": float64(1), "b": "x"},
//			unmarshal: func(data []byte, v interface{}) error {
//				data = data[:len(data)-1]
//				return json.Unmarshal(data, v)
//			},
//			marshal: func(v interface{}) ([]byte, error) {
//				bytes, err := json.Marshal(v)
//				bytes = append(bytes, 's')
//				return bytes, err
//			}}},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			SetMarshal(tt.args.marshal)
//			SetUnmarshal(tt.args.unmarshal)
//			gotBytes, err := marshal(tt.args.v)
//			if err != nil {
//				panic(err)
//			}
//			t.Log(gotBytes, string(gotBytes))
//			got := make(map[string]interface{})
//			if err := unmarshal(gotBytes, &got); err != nil {
//				panic(err)
//			}
//			if !reflect.DeepEqual(got, tt.args.v) {
//				t.Errorf("SetUnmarshalAndSetMarshal() got = %v, want %v", got, tt.args.v)
//			}
//		})
//	}
//}
