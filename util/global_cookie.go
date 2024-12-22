package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

type GlobalCookie struct {
	allCookies map[string]map[string]string
	userCookie map[string]string
	mu         sync.Mutex
}

var globalCookieInstance *GlobalCookie
var once sync.Once

func init() {
	GetGlobalCookie()
}

func GetGlobalCookie() *GlobalCookie {
	once.Do(func() {
		globalCookieInstance = &GlobalCookie{
			allCookies: make(map[string]map[string]string),
			userCookie: make(map[string]string),
		}
		globalCookieInstance.loadCookies()
	})
	return globalCookieInstance
}

func (gc *GlobalCookie) loadCookies() {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	allCookiesData, err := ioutil.ReadFile("data/allCookies.json")
	if err == nil {
		json.Unmarshal(allCookiesData, &gc.allCookies)
	} else {
		log.Println("Failed to load allCookies:", err)
	}

	userCookieData, err := ioutil.ReadFile("data/cookie.json")
	if err == nil {
		json.Unmarshal(userCookieData, &gc.userCookie)
	} else {
		log.Println("Failed to load userCookie:", err)
	}
}

func (gc *GlobalCookie) AllCookies() map[string]map[string]string {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	return gc.allCookies
}

func (gc *GlobalCookie) UserCookie() map[string]string {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	return gc.userCookie
}

func (gc *GlobalCookie) UpdateAllCookies(v map[string]map[string]string) {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	gc.allCookies = v
	data, _ := json.Marshal(v)
	ioutil.WriteFile("data/allCookies.json", data, 0644)
}

func (gc *GlobalCookie) UpdateUserCookie(v map[string]string) {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	gc.userCookie = v
	data, _ := json.Marshal(v)
	ioutil.WriteFile("data/cookie.json", data, 0644)
}
