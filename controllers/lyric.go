package controllers

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"qq-music-api/util"
	"strconv"
	"time"
)

func GetLyric(c *gin.Context) {
	songmid := c.Query("songmid")
	raw := c.Query("raw")

	if songmid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "songmid 不能为空",
		})
		return
	}

	apiUrl := "http://c.y.qq.com/lyric/fcgi-bin/fcg_query_lyric_new.fcg"
	data := url.Values{}
	data.Set("songmid", songmid)
	data.Set("pcachetime", strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10))
	data.Set("g_tk", "5381")
	data.Set("loginUin", "0")
	data.Set("hostUin", "0")
	data.Set("inCharset", "utf8")
	data.Set("outCharset", "utf-8")
	data.Set("notice", "0")
	data.Set("platform", "yqq")
	data.Set("needNewCode", "0")

	result, err := util.MakeRequest(apiUrl, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "Failed to make request",
		})
		return
	}

	// Decode the base64 encoded lyrics
	lyric, err := base64.StdEncoding.DecodeString(result["lyric"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "Failed to decode lyric",
		})
		return
	}
	result["lyric"] = string(lyric)

	trans, err := base64.StdEncoding.DecodeString(result["trans"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "Failed to decode trans",
		})
		return
	}
	result["trans"] = string(trans)

	if raw == "1" {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result": 100,
			"data":   result,
		})
	}
}
