// controllers/singer.go
package controllers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"net/http"
	"qq-music-api/util"
	"reflect"
	"strconv"
	"strings"
)

func GetSingerDesc(c *gin.Context) {
	singermid := c.Query("singermid")
	raw := c.Query("raw")

	if singermid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "singermid 不能为空",
		})
		return
	}

	cacheKey := "singer_desc_" + singermid + "_" + raw
	cacheData, found := util.GetCache(cacheKey)
	if found {
		c.JSON(http.StatusOK, cacheData)
		return
	}

	data := map[string]interface{}{
		"singermid":  singermid,
		"format":     "xml",
		"utf8":       1,
		"outCharset": "utf-8",
	}

	result, err := util.MakeRequest("http://c.y.qq.com/splcloud/fcgi-bin/fcg_get_singer_desc.fcg", data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "Failed to make request",
		})
		return
	}

	page, err := util.MakeRequestRaw("https://y.qq.com/n/yqq/singer/"+singermid+".html", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "Failed to make request",
		})
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "Failed to parse HTML",
		})
		return
	}

	info := result["result"].(map[string]interface{})["data"].(map[string]interface{})["info"].(map[string]interface{})
	info["singername"] = doc.Find(".data__name .data__name_txt").Text()

	for _, k := range []string{"basic", "other"} {
		if item, ok := info[k].(map[string]interface{})["item"]; ok && !isArray(item) {
			info[k].(map[string]interface{})["item"] = []interface{}{item}
		}
	}

	if raw != "1" {
		result = map[string]interface{}{
			"result": 100,
			"data":   info,
		}
	}

	c.JSON(http.StatusOK, result)
	util.SetCache(cacheKey, result, 24*60)
}

func isArray(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.Slice
}

func GetSingerAlbum(c *gin.Context) {
	singermid := c.Query("singermid")
	pageNo := c.DefaultQuery("pageNo", "1")
	pageSize := c.DefaultQuery("pageSize", "20")
	raw := c.Query("raw")

	if singermid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "singermid 不能为空",
		})
		return
	}

	cacheKey := "singer_album_" + singermid + "_" + pageNo + "_" + pageSize + "_" + raw
	cacheData, found := util.GetCache(cacheKey)
	if found {
		c.JSON(http.StatusOK, cacheData)
		return
	}

	pn, _ := strconv.Atoi(pageNo)
	ps, _ := strconv.Atoi(pageSize)

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"comm": map[string]interface{}{
				"ct": 24,
				"cv": 0,
			},
			"singerAlbum": map[string]interface{}{
				"method": "get_singer_album",
				"param": map[string]interface{}{
					"singermid": singermid,
					"order":     "time",
					"begin":     (pn - 1) * ps,
					"num":       pageSize,
					"exstatus":  1,
				},
				"module": "music.web_singer_info_svr",
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
		albumData := result["singerAlbum"].(map[string]interface{})["data"].(map[string]interface{})
		cacheData = gin.H{
			"result": 100,
			"data": map[string]interface{}{
				"list":      albumData["list"],
				"id":        albumData["singer_id"],
				"singermid": albumData["singer_mid"],
				"name":      albumData["singer_name"],
				"total":     albumData["total"],
				"pageNo":    pageNo,
				"pageSize":  pageSize,
			},
		}
		c.JSON(http.StatusOK, cacheData)
		util.SetCache(cacheKey, cacheData, 2*60)
	}
}

func GetSingerSongs(c *gin.Context) {
	singermid := c.Query("singermid")
	num := c.DefaultQuery("num", "20")
	raw := c.Query("raw")
	page := c.DefaultQuery("page", "1")

	pn, _ := strconv.Atoi(page)
	ps, _ := strconv.Atoi(num)

	if singermid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "singermid 不能为空",
		})
		return
	}

	cacheKey := "singer_album_" + singermid + "_" + num + "_" + raw + "_" + page
	cacheData, found := util.GetCache(cacheKey)
	if found {
		c.JSON(http.StatusOK, cacheData)
		return
	}

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"comm": map[string]interface{}{
				"ct": 24,
				"cv": 0,
			},
			"singer": map[string]interface{}{
				"method": "get_singer_detail_info",
				"param": map[string]interface{}{
					"sort":      5,
					"singermid": singermid,
					"sin":       (pn - 1) * ps,
					"num":       ps,
				},
				"module": "music.web_singer_info_svr",
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
		resultData := result["singer"].(map[string]interface{})["data"].(map[string]interface{})
		list := resultData["songlist"].([]interface{})
		extras := resultData["extras"].([]interface{})

		for i, o := range list {
			if i < len(extras) {
				for k, v := range extras[i].(map[string]interface{}) {
					o.(map[string]interface{})[k] = v
				}
			}
		}

		cacheData = gin.H{
			"result": 100,
			"data": map[string]interface{}{
				"list":      list,
				"singer":    resultData["singer_info"],
				"desc":      resultData["desc"],
				"total":     resultData["total_song"],
				"num":       ps,
				"singermid": singermid,
			},
		}
		c.JSON(http.StatusOK, cacheData)
		util.SetCache(cacheKey, cacheData, 2*60)
	}
}
