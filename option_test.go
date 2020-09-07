package requests

import (
	"testing"
)

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
	}{
		{name: "1", args: args{method: "GET", url: BaseUrl, options: []Option{Params{"a": "1"}}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := session.Request(tt.args.method, tt.args.url, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request() error = %v, wantErr %v", err, tt.wantErr)
			}
			resp := &testResp{}
			if err := response.Json(resp); (err != nil) != tt.wantErr {
				t.Errorf("response.Json() error = %v, wantErr %v", err, tt.wantErr)
			}
			if resp.Code != 0 {
				t.Errorf("response.Code != 0, response %v", resp)
			}
			t.Log(resp)
		})
	}
}
