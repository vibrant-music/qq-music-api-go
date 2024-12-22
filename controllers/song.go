package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"qq-music-api/util"
	"strings"
)

func GetSongDetail(c *gin.Context) {
	songmid := c.Query("songmid")
	raw := c.Query("raw")

	if songmid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "songmid 不能为空",
		})
		return
	}

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"songinfo": map[string]interface{}{
				"method": "get_song_detail_yqq",
				"module": "music.pf_song_detail_svr",
				"param": map[string]interface{}{
					"song_mid": songmid,
				},
			},
		},
	}

	result, err := util.MakeRequest("http://u.y.qq.com/cgi-bin/musicu.fcg", data)
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
		c.JSON(http.StatusOK, gin.H{
			"result": 100,
			"data":   result["songinfo"].(map[string]interface{})["data"],
		})
	}
}

func GetSongDownloadURL(c *gin.Context) {
	obj := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		obj[k] = v[0]
	}
	for k, v := range c.Request.PostForm {
		obj[k] = v[0]
	}
	globalCookie := util.GetGlobalCookie()
	userCookie := globalCookie.UserCookie()
	uin := userCookie["uin"]
	qqmusicKey := userCookie["qqmusic_key"]

	id := obj["id"]
	if id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "id 不能为空",
		})
		return
	}

	typeMap := map[string]map[string]string{
		"m4a":  {"s": "C400", "e": ".m4a"},
		"128":  {"s": "M500", "e": ".mp3"},
		"320":  {"s": "M800", "e": ".mp3"},
		"ape":  {"s": "A000", "e": ".ape"},
		"flac": {"s": "F000", "e": ".flac"},
	}
	typeObj, ok := typeMap[obj["type"]]
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "type 传错了，看看文档去",
		})
		return
	}

	file := typeObj["s"] + id + obj["mediaId"] + typeObj["e"]
	guid := fmt.Sprintf("%d", rand.Intn(10000000))

	cacheKey := "song_url_" + file
	cacheData, found := util.GetCache(cacheKey)
	if found {
		c.JSON(http.StatusOK, cacheData)
		return
	}

	var purl, domain string
	for count := 0; purl == "" && count < 10; count++ {
		data := map[string]interface{}{
			"-":           "getplaysongvkey",
			"g_tk":        5381,
			"loginUin":    uin,
			"hostUin":     0,
			"format":      "json",
			"inCharset":   "utf8",
			"outCharset":  "utf-8",
			"platform":    "yqq.json",
			"needNewCode": 0,
			"data": map[string]interface{}{
				"req_0": map[string]interface{}{
					"module": "vkey.GetVkeyServer",
					"method": "CgiGetVkey",
					"param": map[string]interface{}{
						"filename":  []string{file},
						"guid":      guid,
						"songmid":   []string{id},
						"songtype":  []int{0},
						"uin":       uin,
						"loginflag": 1,
						"platform":  "20",
					},
				},
				"comm": map[string]interface{}{
					"uin":    userCookie,
					"format": "json",
					"ct":     19,
					"cv":     0,
					"authst": qqmusicKey,
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

		if result["req_0"].(map[string]interface{})["data"] == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": 400,
				"errMsg": "获取链接出错，建议检查是否携带 cookie",
			})
			return
		}

		midurlinfo := result["req_0"].(map[string]interface{})["data"].(map[string]interface{})["midurlinfo"].([]interface{})
		if len(midurlinfo) > 0 {
			purl = midurlinfo[0].(map[string]interface{})["purl"].(string)
		}

		if domain == "" {
			sip := result["req_0"].(map[string]interface{})["data"].(map[string]interface{})["sip"].([]interface{})
			for _, s := range sip {
				if !strings.HasPrefix(s.(string), "http://ws") {
					domain = s.(string)
					break
				}
			}
			if domain == "" {
				domain = sip[0].(string)
			}
		}
	}

	if purl == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 400,
			"errMsg": "获取播放链接出错",
		})
		return
	}

	if obj["isRedirect"] == "1" {
		c.Redirect(http.StatusFound, domain+purl)
		return
	}

	cacheData = gin.H{
		"data":   domain + purl,
		"result": 100,
	}
	util.SetCache(cacheKey, cacheData, 24*3600)

	c.JSON(http.StatusOK, cacheData)
}

func GetSongPlayURL(c *gin.Context) {
	obj := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		obj[k] = v[0]
	}
	for k, v := range c.Request.PostForm {
		obj[k] = v[0]
	}

	globalCookie := util.GetGlobalCookie()
	userCookie := globalCookie.UserCookie()
	uin := userCookie["uin"]

	id := obj["id"]
	if id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "id 不能为空",
		})
		return
	}

	idArr := strings.Split(id, ",")
	idStr := strings.Join(idArr, `","`)
	cacheKey := "song_url_" + idStr
	cacheData, found := util.GetCache(cacheKey)
	if found {
		c.JSON(http.StatusOK, cacheData)
		return
	}

	url := fmt.Sprintf(`https://u.y.qq.com/cgi-bin/musicu.fcg?-=getplaysongvkey&g_tk=5381&loginUin=%s&hostUin=0&format=json&inCharset=utf8&outCharset=utf-8&notice=0&platform=yqq.json&needNewCode=0&data=%%7B"req_0"%%3A%%7B"module"%%3A"vkey.GetVkeyServer"%%2C"method"%%3A"CgiGetVkey"%%2C"param"%%3A%%7B"guid"%%3A"2796982635"%%2C"songmid"%%3A%%5B"%s"%%5D%%2C"songtype"%%3A%%5B0%%5D%%2C"uin"%%3A"%s"%%2C"loginflag"%%3A1%%2C"platform"%%3A"20"%%7D%%7D%%2C"comm"%%3A%%7B"uin"%%3A%s%%2C"format"%%3A"json"%%2C"ct"%%3A24%%2C"cv"%%3A0%%7D%%7D`, uin, idStr, uin, uin)

	var result map[string]interface{}
	var err error
	for count := 0; count < 5; count++ {
		result, err = util.MakeRequest(url, nil)
		if err == nil && result["req_0"].(map[string]interface{})["data"].(map[string]interface{})["testfile2g"] != nil {
			break
		}
	}

	if err != nil || result["req_0"] == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 200,
			"errMsg": "获取链接失败，建议检查是否登录",
		})
		return
	}

	sip := result["req_0"].(map[string]interface{})["data"].(map[string]interface{})["sip"].([]interface{})
	domain := ""
	for _, s := range sip {
		if !strings.HasPrefix(s.(string), "http://ws") {
			domain = s.(string)
			break
		}
	}
	if domain == "" {
		domain = sip[0].(string)
	}

	data := make(map[string]string)
	midurlinfo := result["req_0"].(map[string]interface{})["data"].(map[string]interface{})["midurlinfo"].([]interface{})
	for _, item := range midurlinfo {
		purl := item.(map[string]interface{})["purl"].(string)
		if purl != "" {
			data[item.(map[string]interface{})["songmid"].(string)] = domain + purl
		}
	}

	cacheData = gin.H{
		"data":   data,
		"result": 100,
	}
	util.SetCache(cacheKey, cacheData, 24*3600)
	c.JSON(http.StatusOK, cacheData)
}
