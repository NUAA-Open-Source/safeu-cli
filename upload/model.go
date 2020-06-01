package upload

import "os"

// 上传相关数据结构

// 上传策略
type UploadPolicy struct {
	AccessID  string
	Host      string
	Expire    int64
	Signature string
	Policy    string
	Dir       string
	Callback  string
}

// 阿里云OSS 接口返回 uuid
type UploadResponse struct {
	UUID string
}

// 上传文件对应状态
const (
	UploadFileReady = iota // 开始生成枚举值, 默认为0
	UploadFileSuccess
	UploadFileFailed
)

// 上传文件
type UploadFile struct {
	File           os.File
	UploadResponse UploadResponse
	StatusCode     int
}

// 上传实例
// 一次上传策略对应一个上传实例
// 一个上传实例包含多个文件上传
type Instance struct {
	UploadPolicy UploadPolicy
	UploadFiles  []UploadFile
}
