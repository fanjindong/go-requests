package requests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"os"
	"testing"
	"time"
)

type testResp struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func (r testResp) String() string {
	return fmt.Sprintf("{Code: %d, Message: %s, Data: %+v }", r.Code, r.Message, r.Data)
}

func getHandler(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": map[string]interface{}{"url": c.Request.URL.String()}})
}

func postHandler(c *gin.Context) {
	reqJson := make(map[string]interface{})

	if p := c.Query("params"); p != "" {
		reqJson["params"] = p
	}

	type form struct {
		Data string `form:"form"`
	}
	f := form{}
	_ = c.ShouldBindWith(&f, binding.FormPost)
	if f.Data != "" {
		reqJson["form"] = f.Data
	}

	if h := c.Request.Header.Get("headers"); h != "" {
		reqJson["headers"] = h
	}
	if cookies := c.Request.Cookies(); len(cookies) > 0 {
		for _, cookie := range cookies {
			reqJson[cookie.Name] = cookie.Value
		}
	}

	_ = c.ShouldBindJSON(&reqJson)
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": reqJson})
}

func timeoutHandler(c *gin.Context) {
	time.Sleep(3 * time.Second)
	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

func TestMain(m *testing.M) {
	r := gin.Default()

	r.GET("/", getHandler)
	r.POST("/", postHandler)
	r.PUT("/", postHandler)
	r.GET("/timeout", timeoutHandler)

	// 监听并在 0.0.0.0:8080 上启动服务
	go func() { _ = r.Run() }()

	session = NewSession()

	code := m.Run()
	os.Exit(code)
}
