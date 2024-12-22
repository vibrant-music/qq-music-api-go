package util

import "encoding/json"

func ConvertToNumber(str string) string {
	// Implement the conversion logic here
	return str
}

func Stringify(data interface{}) string {
	bytes, _ := json.Marshal(data)
	return string(bytes)
}
