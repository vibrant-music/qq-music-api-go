package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"reflect"
)

func GetSign(obj interface{}) (string, error) {
	var str string
	switch v := obj.(type) {
	case string:
		str = v
	case map[string]interface{}:
		bytes, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		str = string(bytes)
	default:
		return "", errors.New("unsupported type")
	}
	return generateSign(str), nil
}

func generateSign(data string) string {
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

var globalObject = func() interface{} {
	if self := reflect.ValueOf("self"); self.IsValid() {
		return self
	}
	if window := reflect.ValueOf("window"); window.IsValid() {
		return window
	}
	if global := reflect.ValueOf("global"); global.IsValid() {
		return global
	}
	panic("unable to locate global object")
}()

func init() {
	globalObject.(map[string]interface{})["__sign_hash_20200305"] = func(e string) string {
		return generateSign(e)
	}
}
