package upload

import (
	"encoding/json"
	"fmt"
	"github.com/arcosx/Nuwa/util"
	"io/ioutil"
	"net/http"
	"strings"
)

// 第 0 步 获取CSRF Token
func (u *Instance) getCSRF() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", util.SAFEU_BASE_URL+"/csrf", nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	u.CSRF = resp.Header.Get("X-Csrf-Token")
	u.Cookie = resp.Header.Get("Set-Cookie")
	return nil
}

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
func (u *Instance) ready(fileFullPaths []string) error {
	for _, fileFullPath := range fileFullPaths {
		var uploadFile UploadFile
		uploadFile.StatusCode = UploadFileReadyCode
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
func (u *Instance) run() []error {
	var errors []error
	for key, file := range u.UploadFiles {
		if file.StatusCode == UploadFileReadyCode {
			err := file.upload()
			if err != nil {
				fmt.Println("upload File failed", err)
				file.StatusCode = UploadFileFailedCode
				errors = append(errors, err)
			} else {
				file.StatusCode = UploadFileSuccessCode
			}
			u.UploadFiles[key] = file
		}
	}
	return errors
}

// 第四步汇总上传结果中的 uuid 发送 finish 函数
func (u *Instance) finish() error {
	var finishRequest FinishRequest
	for _, file := range u.UploadFiles {
		if file.StatusCode == UploadFileSuccessCode {
			finishRequest.Files = append(finishRequest.Files, file.UploadResponse.UUID)
		}
	}
	jsonStr, err := json.Marshal(finishRequest)
	if err != nil {
		fmt.Println("finish json marshal error", err)
	}

	resp, err := requestFinish(string(jsonStr), u.CSRF, u.Cookie)

	if err != nil {
		fmt.Println("finish request post failed", err)
		return err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("finish request reponse not 200")
		return fmt.Errorf("finish request reponse return code: %d ,content %s", resp.StatusCode, respBody)
	}

	// 回填auth token 以及 提取码
	var finishResponse FinishResponse
	err = json.Unmarshal(respBody, &finishResponse)
	if err != nil {
		fmt.Println("finish json unmarshal failed", err)
		return err
	}
	u.Owner = finishResponse.Owner
	u.Recode = finishResponse.Recode
	return nil
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

func requestFinish(body string, csrfToken string, cookie string) (response *http.Response, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", util.SAFEU_BASE_URL, "/v1/upload/finish"), strings.NewReader(body))

	if err != nil {
		fmt.Println("finish request create NewRequest failed", err)
		return response, err
	}
	req.Header.Set("x-csrf-token", csrfToken)
	req.Header.Set("cookie", cookie)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("origin", "https://safeu.a2os.club")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("finish request call failed", err)
		return response, err
	}
	return resp, nil
}
