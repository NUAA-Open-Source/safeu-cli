package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/arcosx/Nuwa/util"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// 获取上传需要的策略参数
func (u *Instance) requestUploadPolicy() (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", util.SAFEU_BASE_URL+"/v1/upload/policy", nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(util.InfoCode["F01"], err)
		return nil, err
	}
	return resp, nil
}

// 解析上传需要的策略参数
func (u *Instance) getUploadPolicy() error {
	var uploadPolicy UploadPolicy
	resp, err := u.requestUploadPolicy()
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

// 构造文件上传请求
// fileFullPath 文件完全路径 /home/just/pig.jpg
func (u *Instance) buildUploadRequest(fileFullPath string) (client *http.Client, url string, values map[string]io.Reader, err error) {
	// TODO: http client 优化及可配置
	client = &http.Client{}
	// fileName:pig.jpg
	fileName := path.Base(fileFullPath)

	file, err := os.Open(fileFullPath)
	if err != nil {
		fmt.Println("buildUploadRequest open file failed ", err, "fileFullPath", fileFullPath)
		return client, url, values, err
	}
	fmt.Println("buildUploadRequest file ", fileName, "open success")
	url = fmt.Sprintf("https://%s", u.UploadPolicy.Host)
	values = map[string]io.Reader{
		"name":                  strings.NewReader(fileName),
		"key":                   strings.NewReader(u.UploadPolicy.Dir + fileName),
		"policy":                strings.NewReader(u.UploadPolicy.Policy),
		"OSSAccessKeyId":        strings.NewReader(u.UploadPolicy.AccessID),
		"success_action_status": strings.NewReader("200"),
		"signature":             strings.NewReader(u.UploadPolicy.Signature),
		"callback":              strings.NewReader(u.UploadPolicy.Callback),
		"file":                  file,
	}
	return client, url, values, nil
}

// 核心函数 上传文件
// TODO: 上传进度条
func (u *Instance) upload(client *http.Client, url string, values map[string]io.Reader) (respBody []byte, err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// 如果是文件标示符 添加文件
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// 添加其他表单信息
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return respBody, err
		}

	}

	_ = w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	uploadBeginTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	// 读取返回响应
	respBody, _ = ioutil.ReadAll(resp.Body)
	uploadUseTime := time.Since(uploadBeginTime)
	fmt.Println("upload use time:", uploadUseTime)

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("upload file to aliyun oss failed, status code:%s,response body: %s", resp.Status, string(respBody))
		return respBody, err
	}
	return
}

// 核心函数 上传文件后的处理
func (u *Instance) handleUploadResponse(respBody []byte) (err error) {
	var uploadResponse UploadResponse
	err = json.Unmarshal(respBody, &uploadResponse)
	if err != nil {
		fmt.Println("handleUploadResponse json unmarshal failed", err)
		return err
	}
	return
}
