// controllers/mv.go
package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qq-music-api/util"
	"strconv"
	"strings"
)

func GetMVInfo(c *gin.Context) {
	id := c.Query("id")
	raw := c.Query("raw")

	if id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "id 不能为空",
		})
		return
	}

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"comm": map[string]interface{}{
				"ct": 24,
				"cv": 4747474,
			},
			"mvinfo": map[string]interface{}{
				"module": "video.VideoDataServer",
				"method": "get_video_info_batch",
				"param": map[string]interface{}{
					"vidlist": []string{id},
					"required": []string{
						"vid", "type", "sid", "cover_pic", "duration", "singers",
						"video_switch", "msg", "name", "desc", "playcnt", "pubdate",
						"isfav", "gmid",
					},
				},
			},
			"other": map[string]interface{}{
				"module": "video.VideoLogicServer",
				"method": "rec_video_byvid",
				"param": map[string]interface{}{
					"vid": id,
					"required": []string{
						"vid", "type", "sid", "cover_pic", "duration", "singers",
						"video_switch", "msg", "name", "desc", "playcnt", "pubdate",
						"isfav", "gmid", "uploader_headurl", "uploader_nick",
						"uploader_encuin", "uploader_uin", "uploader_hasfollow",
						"uploader_follower_num",
					},
					"support": 1,
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
		c.JSON(http.StatusOK, gin.H{
			"result": 100,
			"data": map[string]interface{}{
				"info":      result["mvinfo"].(map[string]interface{})["data"].(map[string]interface{})[id],
				"recommend": result["other"].(map[string]interface{})["data"].(map[string]interface{})["list"],
			},
		})
	}
}

func GetMVUrl(c *gin.Context) {
	id := c.Query("id")
	raw := c.Query("raw")

	if id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "id 不能为空",
		})
		return
	}

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"getMvUrl": map[string]interface{}{
				"module": "gosrf.Stream.MvUrlProxy",
				"method": "GetMvUrls",
				"param": map[string]interface{}{
					"vids":          strings.Split(id, ","),
					"request_typet": 10001,
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
		return
	}

	mvData := result["getMvUrl"].(map[string]interface{})["data"].(map[string]interface{})
	for vid, mv := range mvData {
		mp4Arr := []string{}
		for _, obj := range mv.(map[string]interface{})["mp4"].([]interface{}) {
			if freeflowUrl := obj.(map[string]interface{})["freeflow_url"]; freeflowUrl != nil {
				urls := freeflowUrl.([]interface{})
				if len(urls) > 0 {
					mp4Arr = append(mp4Arr, urls[len(urls)-1].(string))
				}
			}
		}
		mvData[vid] = mp4Arr
	}

	c.JSON(http.StatusOK, gin.H{
		"result": 100,
		"data":   mvData,
	})
}

func GetMVCategory(c *gin.Context) {
	raw := c.Query("raw")

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"comm": map[string]interface{}{
				"ct": 24,
			},
			"mv_tag": map[string]interface{}{
				"module": "MvService.MvInfoProServer",
				"method": "GetAllocTag",
				"param":  map[string]interface{}{},
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
		c.JSON(http.StatusOK, gin.H{
			"result": 100,
			"data":   result["mv_tag"].(map[string]interface{})["data"],
		})
	}
}

func GetMVList(c *gin.Context) {
	raw := c.Query("raw")
	pageNo := c.DefaultQuery("pageNo", "1")
	pageSize := c.DefaultQuery("pageSize", "20")
	version := c.DefaultQuery("version", "7")
	area := c.DefaultQuery("area", "15")

	pn, _ := strconv.Atoi(pageNo)
	ps, _ := strconv.Atoi(pageSize)
	data := map[string]interface{}{
		"data": map[string]interface{}{
			"comm": map[string]interface{}{
				"ct": 24,
			},
			"mv_list": map[string]interface{}{
				"module": "MvService.MvInfoProServer",
				"method": "GetAllocMvInfo",
				"param": map[string]interface{}{
					"start":      (pn - 1) * ps,
					"size":       pageSize,
					"version_id": version,
					"area_id":    area,
					"order":      1,
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
		mvListData := result["mv_list"].(map[string]interface{})["data"].(map[string]interface{})
		list := mvListData["list"]
		total := mvListData["total"]
		c.JSON(http.StatusOK, gin.H{
			"result": 100,
			"data": map[string]interface{}{
				"list":     list,
				"total":    total,
				"area":     area,
				"version":  version,
				"pageNo":   pageNo,
				"pageSize": pageSize,
			},
		})
	}
}

func LikeMV(c *gin.Context) {
	id := c.PostForm("id")
	raw := c.PostForm("raw")
	typeParam := c.DefaultPostForm("type", "1")

	if id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "id 不能为空",
		})
		return
	}

	data := map[string]interface{}{
		"uin":               c.GetHeader("Cookie"),
		"g_tk":              1157392233,
		"format":            "json",
		"inCharset":         "utf-8",
		"outCharset":        "utf-8",
		"cmdtype":           typeParam != "1",
		"reqtype":           1,
		"mvidlist":          id,
		"mvidtype":          0,
		"cv":                4747474,
		"ct":                24,
		"notice":            0,
		"platform":          "yqq.json",
		"needNewCode":       1,
		"g_tk_new_20200303": 1859542818,
		"cid":               205361448,
	}

	result, err := util.MakeRequest("https://c.y.qq.com/mv/fcgi-bin/fcg_add_del_myfav_mv.fcg", data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 500,
			"errMsg": "Failed to make request",
		})
		return
	}

	if raw == "1" {
		c.JSON(http.StatusOK, result)
		return
	}

	if result["code"].(float64) != 0 {
		c.JSON(http.StatusOK, gin.H{
			"result": 200,
			"errMsg": result["msg"],
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": 100,
		"data":   "操作成功！",
	})
}
