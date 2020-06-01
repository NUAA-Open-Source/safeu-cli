package util

// 行为码定义
// 成功 0
// 失败 1
var InfoCode = map[string]string{
	"S00": "ALL SUCCESS",
	"S01": "Request Safeu Policy Success",
	"F01": "Request Safeu Policy failed",
	"S02": "Read Upload File Success",
	"F02": "Read Upload File failed",
	"S03": "Upload AliYun OSS Success",
	"F03": "Upload AliYun OSS failed",
}
