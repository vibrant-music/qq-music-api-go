package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"qq-music-api/constant"
	"strings"
)

func MakeRequest(url string, data interface{}) (map[string]interface{}, error) {
	var req *http.Request
	var err error

	if data != nil {
		jsonData, _ := json.Marshal(data)
		req, err = http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	} else {
		req, err = http.NewRequest("GET", url, nil)
	}

	if err != nil {
		return nil, err
	}

	cookieObj := globalCookieInstance.UserCookie()

	var cookies []string
	for k, v := range cookieObj {
		cookies = append(cookies, fmt.Sprintf("%s=%s", k, v))
	}
	req.Header.Set("Cookie", strings.Join(cookies, "; "))
	req.Header.Set("Referer", "https://y.qq.com")
	req.Header.Set("X-XSRF-TOKEN", "XSRF-TOKEN")
	req.Header.Set("withCredentials", "true")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if responseStr := string(body); responseStr != "" {
		responseStr = strings.ReplaceAll(responseStr, "callback(", "")
		responseStr = strings.ReplaceAll(responseStr, "MusicJsonCallback(", "")
		responseStr = strings.ReplaceAll(responseStr, "jsonCallback(", "")
		responseStr = strings.TrimSuffix(responseStr, ")")
		body = []byte(responseStr)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func MakeRequestV2(method constant.HTTPMethod, baseUrl string, data interface{}) (map[string]interface{}, error) {
	var req *http.Request
	var err error

	var buffer *bytes.Buffer
	if data != nil {
		jsonData, _ := json.Marshal(data)
		buffer = bytes.NewBuffer(jsonData)
	}
	switch method {
	case constant.HTTPGet:
		urlWithParams, uErr := url.Parse(baseUrl)
		if uErr != nil {
			return nil, uErr
		}
		query := urlWithParams.Query()
		for key, value := range data.(map[string]interface{}) {
			query.Set(key, fmt.Sprintf("%v", value))
		}
		urlWithParams.RawQuery = query.Encode()
		req, err = http.NewRequest("GET", urlWithParams.String(), nil)
	case constant.HTTPPost:
		req, err = http.NewRequest("POST", baseUrl, buffer)
	case constant.HTTPPut:
		req, err = http.NewRequest("PUT", baseUrl, buffer)
	case constant.HTTPDelete:
		req, err = http.NewRequest("DELETE", baseUrl, buffer)
	case constant.HTTPPatch:
		req, err = http.NewRequest("PATCH", baseUrl, buffer)
	}

	if err != nil {
		return nil, err
	}

	cookieObj := globalCookieInstance.UserCookie()

	var cookies []string
	for k, v := range cookieObj {
		cookies = append(cookies, fmt.Sprintf("%s=%s", k, v))
	}
	req.Header.Set("Cookie", strings.Join(cookies, "; "))
	req.Header.Set("Referer", "https://y.qq.com")
	req.Header.Set("X-XSRF-TOKEN", "XSRF-TOKEN")
	req.Header.Set("withCredentials", "true")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if responseStr := string(body); responseStr != "" {
		responseStr = strings.ReplaceAll(responseStr, "callback(", "")
		responseStr = strings.ReplaceAll(responseStr, "MusicJsonCallback(", "")
		responseStr = strings.ReplaceAll(responseStr, "jsonCallback(", "")
		responseStr = strings.TrimSuffix(responseStr, ")")
		body = []byte(responseStr)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func MakeRequestRaw(url string, data interface{}) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	cookieObj := globalCookieInstance.UserCookie()

	var cookies []string
	for k, v := range cookieObj {
		cookies = append(cookies, fmt.Sprintf("%s=%s", k, v))
	}
	req.Header.Set("Cookie", strings.Join(cookies, "; "))
	req.Header.Set("Referer", "https://y.qq.com")
	req.Header.Set("X-XSRF-TOKEN", "XSRF-TOKEN")
	req.Header.Set("withCredentials", "true")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
