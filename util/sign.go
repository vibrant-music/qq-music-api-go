package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
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
