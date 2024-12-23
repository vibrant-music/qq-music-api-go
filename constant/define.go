package constant

type HTTPMethod string

// 定义常量
const (
	HTTPGet    HTTPMethod = "GET"
	HTTPPost   HTTPMethod = "POST"
	HTTPPut    HTTPMethod = "PUT"
	HTTPDelete HTTPMethod = "DELETE"
	HTTPPatch  HTTPMethod = "PATCH"
)
