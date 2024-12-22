package controllers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"net/http"
	"qq-music-api/util"
	"strings"
)

func GetRecommendPlaylist(c *gin.Context) {
	raw := c.Query("raw")
	pageNo := c.DefaultQuery("pageNo", "1")
	pageSize := c.DefaultQuery("pageSize", "20")
	id := c.DefaultQuery("id", "3317") // 3317: 官方歌单，59：经典，71：情歌，3056：网络歌曲，64：KTV热歌

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"comm": map[string]interface{}{
				"ct": 24,
			},
			"playlist": map[string]interface{}{
				"method": "get_playlist_by_category",
				"param": map[string]interface{}{
					"id":      id,
					"curPage": pageNo,
					"size":    pageSize,
					"order":   5,
					"titleid": id,
				},
				"module": "playlist.PlayListPlazaServer",
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
		playlistData := result["playlist"].(map[string]interface{})["data"].(map[string]interface{})
		total := playlistData["total"]
		v_playlist := playlistData["v_playlist"]
		c.JSON(http.StatusOK, gin.H{
			"result": 100,
			"data": map[string]interface{}{
				"total":    total,
				"list":     v_playlist,
				"id":       id,
				"pageNo":   pageNo,
				"pageSize": pageSize,
			},
		})
	}
}

func GetRecommendPlaylistByUser(c *gin.Context) {
	raw := c.Query("raw")

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"comm": map[string]interface{}{
				"ct": 24,
			},
			"recomPlaylist": map[string]interface{}{
				"method": "get_hot_recommend",
				"param": map[string]interface{}{
					"async": 1,
					"cmd":   2,
				},
				"module": "playlist.HotRecommendServer",
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
		list := result["recomPlaylist"].(map[string]interface{})["data"].(map[string]interface{})["v_hot"]
		c.JSON(http.StatusOK, gin.H{
			"result": 100,
			"data": map[string]interface{}{
				"list":  list,
				"count": len(list.([]interface{})),
			},
		})
	}
}

func GetRecommendDaily(c *gin.Context) {
	c.DefaultQuery("ownCookie", "1")
	page, err := util.MakeRequestRaw("https://c.y.qq.com/node/musicmac/v6/index.html", nil)
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

	firstList := doc.Find(".mod_for_u .playlist__item").First()
	var id string
	if firstList.Find(".playlist__name").Text() == "今日私享" {
		id, _ = firstList.Find(".playlist__link").Attr("data-rid")
	}

	if id == "" {
		c.JSON(http.StatusMovedPermanently, gin.H{
			"result": 301,
			"errMsg": "未登录",
		})
		return
	}

	// todo
}

func GetRecommendBanner(c *gin.Context) {
	c.DefaultQuery("ownCookie", "1")
	page, err := util.MakeRequestRaw("https://c.y.qq.com/node/musicmac/v6/index.html", nil)
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

	var result []map[string]interface{}
	doc.Find(".focus__box .focus__pic").Each(func(i int, s *goquery.Selection) {
		domA := s.Find("a")
		domImg := s.Find("img")
		typeAttr, _ := domA.Attr("data-type")
		idAttr, _ := domA.Attr("data-rid")
		obj := map[string]interface{}{
			"type":   typeAttr,
			"id":     idAttr,
			"picUrl": domImg.AttrOr("src", ""),
			"h5Url": map[string]string{
				"10002": fmt.Sprintf("https://y.qq.com/musicmac/v6/album/detail.html?albumid=%s", idAttr),
			}[typeAttr],
			"typeStr": map[string]string{
				"10002": "album",
			}[typeAttr],
		}
		result = append(result, obj)
	})

	c.JSON(http.StatusOK, gin.H{
		"result": 100,
		"data":   result,
	})
}
