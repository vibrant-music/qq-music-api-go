package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qq-music-api/util"
)

func Search(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"result": 500, "errMsg": "关键词不能为空"})
		return
	}

	// Example URL for search request
	url := "https://example.com/search?key=" + key
	data, err := util.MakeRequest(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": 500, "errMsg": "搜索请求失败"})
		return
	}

	// Process data and return response
	c.JSON(http.StatusOK, gin.H{"result": 100, "data": string(data)})
}

func HotSearch(c *gin.Context) {
	// Example URL for hot search request
	url := "https://example.com/hotsearch"
	data, err := util.MakeRequest(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": 500, "errMsg": "热搜请求失败"})
		return
	}

	// Process data and return response
	c.JSON(http.StatusOK, gin.H{"result": 100, "data": string(data)})
}

func QuickSearch(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"result": 500, "errMsg": "key 不能为空"})
		return
	}

	// Example URL for quick search request
	url := "https://example.com/quicksearch?key=" + key
	data, err := util.MakeRequest(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": 500, "errMsg": "快速搜索请求失败"})
		return
	}

	// Process data and return response
	c.JSON(http.StatusOK, gin.H{"result": 100, "data": string(data)})
}
