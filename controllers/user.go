package controllers

import (
	"errors"
	"net/http"
	"net/url"
	"qq-music-api/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

var c = cache.New(5*time.Minute, 10*time.Minute)

func Cookie(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"result": 100,
		"data":   c.Request.Cookies(),
	})
}

func SetCookie(c *gin.Context) {
	var requestBody struct {
		Data string `json:"data"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": 500, "errMsg": "Invalid request body"})
		return
	}

	data := requestBody.Data
	cookiePairs := strings.Split(data, "; ")
	userCookie := make(map[string]string)
	for _, pair := range cookiePairs {
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			userCookie[parts[0]] = parts[1]
		}
	}

	if userCookie["login_type"] == "2" {
		userCookie["uin"] = userCookie["wxuin"]
	}
	userCookie["uin"] = strings.Trim(userCookie["uin"], " ")

	globalCookie := util.GetGlobalCookie()
	allCookies := globalCookie.AllCookies()
	allCookies[userCookie["uin"]] = userCookie
	globalCookie.UpdateAllCookies(allCookies)
	globalCookie.UpdateUserCookie(userCookie)

	c.Header("Access-Control-Allow-Origin", "https://y.qq.com")
	c.Header("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.JSON(http.StatusOK, gin.H{"result": 100, "data": "操作成功"})
}

func GetCookie(c *gin.Context) {
	id := c.Query("id")
	globalCookie := util.GetGlobalCookie()
	if id == "" {
		cookieObj := globalCookie.UserCookie()
		for k, v := range cookieObj {
			if len(v) < 255 {
				c.SetCookie(k, v, 86400, "/", "", false, true)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"result":  100,
			"message": "获取 cookie 成功",
		})
		return
	}

	cookieObj := globalCookie.AllCookies()[id]
	if cookieObj == nil {
		cookieObj = make(map[string]string)
	}

	for k, v := range cookieObj {
		if len(v) < 255 {
			c.SetCookie(k, v, 86400, "/", "", false, true)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"result":  100,
		"message": "获取 cookie 成功",
	})
}

type QQLoginRequest struct {
	Req1 struct {
		Module string `json:"module"`
		Method string `json:"method"`
		Param  struct {
			ExpiredIn int    `json:"expired_in"`
			MusicID   string `json:"musicid"`
			MusicKey  string `json:"musickey"`
		} `json:"param"`
	} `json:"req1"`
}

func Refresh(c *gin.Context) {
	uin, err := c.Cookie("uin")
	qmKeyst, err1 := c.Cookie("qm_keyst")
	qqmusicKey, err2 := c.Cookie("qqmusic_key")

	if err != nil || (err1 != nil && err2 != nil) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": 301,
			"errMsg": "未登陆",
		})
		return
	}

	musicKey := qmKeyst
	if musicKey == "" {
		musicKey = qqmusicKey
	}

	data := QQLoginRequest{}
	data.Req1.Module = "QQConnectLogin.LoginServer"
	data.Req1.Method = "QQLogin"
	data.Req1.Param.ExpiredIn = 7776000
	data.Req1.Param.MusicID = util.ConvertToNumber(uin)
	data.Req1.Param.MusicKey = musicKey

	sign, _ := util.GetSign(data)
	url := "https://u6.y.qq.com/cgi-bin/musics.fcg?sign=" + sign + "&format=json&inCharset=utf8&outCharset=utf-8&data=" + url.QueryEscape(util.Stringify(data))

	result, err := util.MakeRequest(url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": 200,
			"errMsg": "刷新失败，建议重新设置cookie",
		})
		return
	}

	if req1, ok := result["req1"].(map[string]interface{}); ok {
		if iData, ok2 := req1["data"].(map[string]interface{}); ok2 {
			if musickey, ok3 := iData["musickey"].(string); ok3 {
				c.SetCookie("qm_keyst", musickey, 86400, "/", "", false, true)
				c.SetCookie("qqmusic_key", musickey, 86400, "/", "", false, true)
				c.JSON(http.StatusOK, gin.H{
					"result": 100,
					"data": gin.H{
						"musickey": musickey,
					},
				})
				return
			}
		}
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"result": 200,
		"errMsg": "刷新失败，建议议重新设置cookie",
	})
}

func RefreshFromCache(uin string) error {
	globalCookie := util.GetGlobalCookie()
	var userCookies map[string]string
	if uin == "" {
		userCookies = globalCookie.UserCookie()
		uin = userCookies["uin"]
	} else {
		userCookies = globalCookie.AllCookies()[uin]
	}

	if len(userCookies) == 0 {
		return errors.New("uin not found")
	}

	qmKeyst := userCookies["qm_keyst"]
	qqmusicKey := userCookies["qqmusic_key"]

	if qmKeyst == "" || qqmusicKey == "" {
		return errors.New("musickey not found")
	}

	musicKey := qmKeyst
	if musicKey == "" {
		musicKey = qqmusicKey
	}

	data := QQLoginRequest{}
	data.Req1.Module = "QQConnectLogin.LoginServer"
	data.Req1.Method = "QQLogin"
	data.Req1.Param.ExpiredIn = 7776000
	data.Req1.Param.MusicID = util.ConvertToNumber(uin)
	data.Req1.Param.MusicKey = musicKey

	sign, _ := util.GetSign(data)
	url := "https://u6.y.qq.com/cgi-bin/musics.fcg?sign=" + sign + "&format=json&inCharset=utf8&outCharset=utf-8&data=" + url.QueryEscape(util.Stringify(data))

	result, err := util.MakeRequest(url, nil)
	if err != nil {
		return err
	}

	if req1, ok := result["req1"].(map[string]interface{}); ok {
		if iData, ok2 := req1["data"].(map[string]interface{}); ok2 {
			if musickey, ok3 := iData["musickey"].(string); ok3 {
				userCookie := make(map[string]string)
				userCookie["uin"] = uin
				userCookie["qm_keyst"] = musickey
				userCookie["qqmusic_key"] = musickey
				allCookies := globalCookie.AllCookies()
				allCookies[uin] = userCookie
				globalCookie.UpdateAllCookies(allCookies)
				globalCookie.UpdateUserCookie(userCookie)
				return nil
			}
		}
	}

	return errors.New("RefreshFromCache failed")

}
