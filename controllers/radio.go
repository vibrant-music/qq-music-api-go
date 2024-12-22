package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qq-music-api/util"
)

func GetRadio(c *gin.Context) {
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
			"songlist": map[string]interface{}{
				"module": "mb_track_radio_svr",
				"method": "get_radio_track",
				"param": map[string]interface{}{
					"id":        id,
					"firstplay": 1,
					"num":       15,
				},
			},
			"radiolist": map[string]interface{}{
				"module": "pf.radiosvr",
				"method": "GetRadiolist",
				"param": map[string]interface{}{
					"ct": "24",
				},
			},
			"comm": map[string]interface{}{
				"ct": 24,
				"cv": 0,
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
			"data":   result["songlist"].(map[string]interface{})["data"],
		})
	}
}

func GetRadioCategory(c *gin.Context) {
	raw := c.Query("raw")

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"songlist": map[string]interface{}{
				"module": "mb_track_radio_svr",
				"method": "get_radio_track",
				"param": map[string]interface{}{
					"id":        99,
					"firstplay": 1,
					"num":       15,
				},
			},
			"radiolist": map[string]interface{}{
				"module": "pf.radiosvr",
				"method": "GetRadiolist",
				"param": map[string]interface{}{
					"ct": "24",
				},
			},
			"comm": map[string]interface{}{
				"ct": 24,
				"cv": 0,
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
			"data":   result["radiolist"].(map[string]interface{})["data"].(map[string]interface{})["radio_list"],
		})
	}
}
