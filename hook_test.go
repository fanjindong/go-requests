package requests

import (
	"fmt"
	"testing"
)

type mockHook struct {
}

func (m mockHook) BeforeProcess(req *Request) {
	fmt.Println("before process")
}

func (m mockHook) AfterProcess(req *Request, resp *Response, err error) {
	fmt.Println("after process")
}

func TestHook(t *testing.T) {
	client := NewClient()
	client.AddHook(mockHook{})
	t.Log(client.Get(testUrl))
}
