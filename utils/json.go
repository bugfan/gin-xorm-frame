package utils

import "github.com/bugfan/jsoniter"

// 通用深层json解析器
/*
*	功能: 加速反序列化 直接读取json
*   参照文档: http://jsoniter.com/migrate-from-go-std.html
 */
func NewIJSON() jsoniter.API {
	return jsoniter.ConfigCompatibleWithStandardLibrary
}
