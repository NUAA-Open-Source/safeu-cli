package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/arcosx/Nuwa/util"
	"io/ioutil"
	"net/http"
)

// 第一步 获取上传需要的策略参数
func (u *Instance) getUploadPolicy() error {
	var uploadPolicy UploadPolicy
	resp, err := requestUploadPolicy()
	if err != nil {
		fmt.Println("getUploadPolicy requestUploadPolicy failed", err)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &uploadPolicy)
	if err != nil {
		fmt.Println("getUploadPolicy json unmarshal failed", err)
		return err
	}
	u.UploadPolicy = uploadPolicy
	return nil
}

// 第二步 上传文件准备
// fileFullPaths 文件完全路径
func (u Instance) ready(fileFullPaths []string) error {
	for _, fileFullPath := range fileFullPaths {
		var uploadFile UploadFile
		uploadFile.statusCode = UploadFileReadyCode
		err := uploadFile.buildUploadRequest(u.UploadPolicy, fileFullPath)
		if err != nil {
			return err
		}
		u.UploadFiles = append(u.UploadFiles, uploadFile)
	}
	return nil
}

// 第三步 开始上传文件
// TODO: 并发
func (u Instance) run() {
	for _, file := range u.UploadFiles {
		if file.statusCode == UploadFileReadyCode {
			err := file.upload()
			if err != nil {
				fmt.Println("upload file failed", err)
				file.statusCode = UploadFileFailedCode
			} else {
				file.statusCode = UploadFileSuccessCode
			}
		}
	}
}

// 第四步汇总上传结果中的 uuid 发送 finish 函数
func (u Instance) finish() error {
	var finishRequest FinishRequest
	for _, file := range u.UploadFiles {
		if file.statusCode == UploadFileSuccessCode {
			finishRequest.Files = append(finishRequest.Files, file.uploadResponse.UUID)
		}
	}
	jsonStr, err := json.Marshal(finishRequest)
	if err != nil {
		fmt.Println("finish json marshal error", err)
	}
	resp, err := http.Post(fmt.Sprintf("%s%s", util.SAFEU_BASE_URL, "/v1/upload/finish"), "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		fmt.Println("finish")
	}
	respBody, _ := ioutil.ReadAll(resp.Body)

}

// 辅助函数

// 获取上传需要的策略参数
func requestUploadPolicy() (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", util.SAFEU_BASE_URL+"/v1/upload/policy", nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(util.InfoCode["F01"], err)
		return nil, err
	}
	return resp, nil
}
