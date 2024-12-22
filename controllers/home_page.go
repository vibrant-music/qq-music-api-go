package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":   "QQ 音乐 api",
		"content": `<a href="http://jsososo.github.io/QQMusicApi">查看文档</a>`,
	})
}
