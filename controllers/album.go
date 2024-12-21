package controllers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"net/http"
	"qq-music-api/util"
	"strings"
)

func AlbumInfo(c *gin.Context) {
	albummid := c.Query("albummid")

	if albummid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"result": 500, "errMsg": "albummid 不能为空"})
		return
	}

	pageInfo, err := util.MakeRequestRaw("https://y.qq.com/n/yqq/album/"+albummid+".html", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": 400, "errMsg": err.Error()})
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(pageInfo))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": 400, "errMsg": err.Error()})
		return
	}

	var albumInfo map[string]interface{}
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		s.Contents().Each(func(i int, content *goquery.Selection) {
			data := content.Text()
			if strings.Contains(data, "window.__USE_SSR__") {
				data = strings.Replace(data, "window.__", "window__", -1)
				albumInfo = util.EvalJS(data)
			}
		})
	})

	if albumInfo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": 400, "errMsg": "Failed to parse album info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": 100, "data": albumInfo})
}

func AlbumSongs(c *gin.Context) {
	albummid := c.Query("albummid")
	raw := c.Query("raw")

	if albummid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"result": 500, "errMsg": "albummid 不能为空"})
		return
	}

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"comm": map[string]interface{}{
				"ct": 24,
				"cv": 10000,
			},
			"albumSonglist": map[string]interface{}{
				"method": "GetAlbumSongList",
				"param": map[string]interface{}{
					"albumMid": albummid,
					"albumID":  0,
					"begin":    0,
					"num":      999,
					"order":    2,
				},
				"module": "music.musichallAlbum.AlbumSongList",
			},
		},
	}

	result, err := util.MakeRequest("https://u.y.qq.com/cgi-bin/musicu.fcg?g_tk=5381&format=json&inCharset=utf8&outCharset=utf-8", data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": 500, "errMsg": "快速搜索请求失败"})
		return
	}

	if raw == "1" {
		c.JSON(http.StatusOK, result)
		return
	}

	resData := gin.H{
		"result": 100,
		"data": gin.H{
			"list":     result["albumSonglist"].(map[string]interface{})["data"].(map[string]interface{})["songList"],
			"total":    result["albumSonglist"].(map[string]interface{})["data"].(map[string]interface{})["totalNum"],
			"albummid": result["albumSonglist"].(map[string]interface{})["data"].(map[string]interface{})["albumMid"],
		},
	}

	c.JSON(http.StatusOK, resData)
}
