package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"qq-music-api/constant"
	"qq-music-api/model/search"
	"qq-music-api/util"
	"strconv"
)

func Search(c *gin.Context) {
	pageNo, _ := strconv.Atoi(c.DefaultQuery("pageNo", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	key := c.Query("key")
	t, _ := strconv.Atoi(c.DefaultQuery("t", "0"))
	raw := c.Query("raw")

	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"result": 500, "errMsg": "关键词不能为空"})
		return
	}

	cacheKey := "search_" + key + "_" + strconv.Itoa(pageNo) + "_" + strconv.Itoa(pageSize) + "_" + strconv.Itoa(t)
	cacheData, found := util.GetCache(cacheKey)
	if found {
		c.JSON(http.StatusOK, cacheData)
		return
	}

	url := map[int]string{
		2: "https://c.y.qq.com/soso/fcgi-bin/client_music_search_songlist?remoteplace=txt.yqq.playlist&page_no=" + strconv.Itoa(pageNo-1) + "&num_per_page=" + strconv.Itoa(pageSize) + "&query=" + key,
	}[t]
	if url == "" {
		url = "http://c.y.qq.com/soso/fcgi-bin/client_search_cp"
	}

	typeMap := map[int]string{
		0:  "song",
		2:  "songlist",
		7:  "lyric",
		8:  "album",
		12: "mv",
		9:  "singer",
	}

	if _, ok := typeMap[t]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"result": 500, "errMsg": "搜索类型错误，检查一下参数 t"})
		return
	}

	data := map[string]interface{}{
		"format": "json",
		"n":      pageSize,
		"p":      pageNo,
		"w":      key,
		"cr":     1,
		"g_tk":   5381,
		"t":      t,
	}

	if t == 2 {
		data = map[string]interface{}{
			"query":        key,
			"page_no":      pageNo - 1,
			"num_per_page": pageSize,
		}
	}

	result, err := util.MakeRequestV2(constant.HTTPGet, url, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": 500, "errMsg": "搜索请求失败"})
		return
	}

	if raw == "1" {
		c.JSON(http.StatusOK, result)
		return
	}

	var rsp search.Response
	resJson, _ := json.Marshal(result)
	err = jsoniter.Unmarshal(resJson, &rsp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": 500, "errMsg": fmt.Sprintf("搜索请求失败 %v", err)})
		return
	}

	// Process and format the result data
	keyword := rsp.Data.Keyword
	//keyMap := map[int]string{
	//	0:  "song",
	//	2:  "",
	//	7:  "lyric",
	//	8:  "album",
	//	12: "mv",
	//	9:  "singer",
	//}
	searchResult := rsp.Data.Song
	list := searchResult.List
	curpage := searchResult.CurPage
	curnum := searchResult.CurNum
	totalnum := searchResult.TotalNum

	var total int
	switch t {
	case 2:
		//pageNo = page_no + 1
		//pageSize = num_per_page
		//total = display_num
	default:
		pageNo = curpage
		pageSize = curnum
		total = totalnum
	}

	resData := gin.H{
		"result": 100,
		"data": gin.H{
			"list":     list,
			"pageNo":   pageNo,
			"pageSize": pageSize,
			"total":    total,
			"key":      keyword,
			"t":        t,
			"type":     typeMap[t],
		},
	}

	util.SetCache(cacheKey, resData, constant.OneDaySeconds)
	c.JSON(http.StatusOK, resData)
}

func HotSearch(c *gin.Context) {
	raw := c.Query("raw")
	url := "https://c.y.qq.com/splcloud/fcgi-bin/gethotkey.fcg"

	result, err := util.MakeRequest(url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": 500, "errMsg": "热搜请求失败"})
		return
	}

	if raw == "1" {
		c.JSON(http.StatusOK, result)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": 100, "data": result["data"].(map[string]interface{})["hotkey"]})
}

func QuickSearch(c *gin.Context) {
	key := c.Query("key")
	raw := c.Query("raw")

	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"result": 500, "errMsg": "key ?"})
		return
	}

	url := "https://c.y.qq.com/splcloud/fcgi-bin/smartbox_new.fcg?key=" + key + "&g_tk=5381"
	result, err := util.MakeRequest(url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": 500, "errMsg": "快速搜索请求失败"})
		return
	}

	if raw == "1" {
		c.JSON(http.StatusOK, result)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": 100, "data": result["data"]})
}
