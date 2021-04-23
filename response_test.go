package requests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponse(t *testing.T) {
	url := baseUrl
	resp, err := Get(url)
	assert.NoError(t, err)

	assert.Equal(t, "application/json; charset=utf-8", resp.Headers.Get("Content-Type"))
	assert.Equal(t, "application/json; charset=utf-8", resp.Headers.Get("content-type"))
}

func TestResponse_SetEncode(t *testing.T) {
	resp := &Response{bytes: []byte("你好")}
	err := resp.SetEncode("utf-8")
	assert.NoError(t, err)
	t.Log(resp.GetEncode())
	t.Log(resp.Text())

	err = resp.SetEncode("GBK")
	assert.NoError(t, err)
	t.Log(resp.GetEncode())
	t.Log(resp.Text())

	err = resp.SetEncode("ASCII")
	assert.NoError(t, err)
	t.Log(resp.GetEncode())
	t.Log(resp.Text())
}
