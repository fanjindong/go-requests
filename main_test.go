package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"os"
	"testing"
)

func getHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"url": c.Request.URL.String(),
	})
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

	_ = c.ShouldBindJSON(&reqJson)
	c.JSON(200, gin.H{
		"data": reqJson,
	})
}

func TestMain(m *testing.M) {
	r := gin.Default()

	r.GET("/get", getHandler)
	r.POST("/post", postHandler)
	r.PUT("/put", postHandler)

	go r.Run() // 监听并在 0.0.0.0:8080 上启动服务

	code := m.Run()
	os.Exit(code)
}
