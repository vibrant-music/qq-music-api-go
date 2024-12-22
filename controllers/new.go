// controllers/music.go
package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qq-music-api/util"
)

func GetNewSongs(c *gin.Context) {
	typeParam := c.DefaultQuery("type", "0")
	raw := c.Query("raw")

	newType := map[string]int{
		"0": 5, // 最新
		"1": 1, // 内地
		"2": 6, // 港台
		"3": 2, // 欧美
		"4": 4, // 韩国
		"5": 3, // 日本
	}[typeParam]

	if newType == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "type 不合法",
		})
		return
	}

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"comm": map[string]interface{}{
				"ct": 24,
			},
			"new_song": map[string]interface{}{
				"module": "newsong.NewSongServer",
				"method": "get_new_song_info",
				"param": map[string]interface{}{
					"type": newType,
				},
			},
		},
	}

	result, err := util.MakeRequest("https://u.y.qq.com/cgi-bin/musicu.fcg", data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "Failed to make request",
		})
		return
	}

	if raw == "1" {
		c.JSON(http.StatusOK, result)
	} else {
		newSongData := result["new_song"].(map[string]interface{})["data"].(map[string]interface{})
		lan := newSongData["lan"]
		songList := newSongData["songlist"]
		c.JSON(http.StatusOK, gin.H{
			"result": 100,
			"data": map[string]interface{}{
				"lan":  lan,
				"list": songList,
				"type": newType,
			},
		})
	}
}

func GetNewAlbum(c *gin.Context) {
	typeParam := c.DefaultQuery("type", "1")
	num := c.DefaultQuery("num", "10")
	raw := c.Query("raw")

	typeName := map[string]string{
		"1": "内地",
		"2": "港台",
		"3": "欧美",
		"4": "韩国",
		"5": "日本",
		"6": "其他",
	}[typeParam]

	if typeName == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "type 不合法",
		})
		return
	}

	data := map[string]interface{}{
		"platform":    "yqq.json",
		"needNewCode": 0,
		"data": map[string]interface{}{
			"comm": map[string]interface{}{
				"ct": 24,
			},
			"new_album": map[string]interface{}{
				"module": "newalbum.NewAlbumServer",
				"method": "get_new_album_info",
				"param": map[string]interface{}{
					"area": typeParam,
					"sin":  0,
					"num":  num,
				},
			},
		},
	}

	result, err := util.MakeRequest("https://u.y.qq.com/cgi-bin/musicu.fcg", data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "Failed to make request",
		})
		return
	}

	if raw == "1" {
		c.JSON(http.StatusOK, result)
	} else {
		newAlbumData := result["new_album"].(map[string]interface{})["data"].(map[string]interface{})
		albums := newAlbumData["albums"]
		c.JSON(http.StatusOK, gin.H{
			"result": 100,
			"data": map[string]interface{}{
				"list":     albums,
				"type":     typeParam,
				"typeName": typeName,
			},
		})
	}
}

func GetNewMV(c *gin.Context) {
	typeParam := c.DefaultQuery("type", "0")
	raw := c.Query("raw")

	typeName := map[string]string{
		"0": "精选",
		"1": "内地",
		"2": "港台",
		"3": "欧美",
		"4": "韩国",
		"5": "日本",
	}[typeParam]

	lan := map[string]string{
		"0": "all",
		"1": "neidi",
		"2": "gangtai",
		"3": "oumei",
		"4": "korea",
		"5": "janpan",
	}[typeParam]

	if lan == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "type 不合法",
		})
		return
	}

	data := map[string]interface{}{
		"cmd": "shoubo",
		"lan": lan,
	}

	result, err := util.MakeRequest("https://c.y.qq.com/mv/fcgi-bin/getmv_by_tag", data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "Failed to make request",
		})
		return
	}

	if raw == "1" {
		c.JSON(http.StatusOK, result)
	} else {
		mvData := result["data"].(map[string]interface{})
		mvlist := mvData["mvlist"]
		c.JSON(http.StatusOK, gin.H{
			"result": 100,
			"data": map[string]interface{}{
				"list":     mvlist,
				"lan":      mvData["lan"],
				"typeName": typeName,
			},
		})
	}
}
