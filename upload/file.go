package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// 构造文件上传请求
// fileFullPath 文件完全路径 /home/just/pig.jpg
func (f *UploadFile) buildUploadRequest(uploadPolicy UploadPolicy, fileFullPath string) (err error) {
	// TODO: http client 优化及可配置
	f.client = &http.Client{}
	// fileName:pig.jpg
	fileName := path.Base(fileFullPath)

	f.file, err = os.Open(fileFullPath)
	if err != nil {
		fmt.Println("buildUploadRequest open file failed ", err, "fileFullPath", fileFullPath)
		return err
	}
	fmt.Println("buildUploadRequest file ", fileName, "open success")
	f.url = fmt.Sprintf("https://%s", uploadPolicy.Host)
	f.values = map[string]io.Reader{
		"name":                  strings.NewReader(fileName),
		"key":                   strings.NewReader(uploadPolicy.Dir + fileName),
		"policy":                strings.NewReader(uploadPolicy.Policy),
		"OSSAccessKeyId":        strings.NewReader(uploadPolicy.AccessID),
		"success_action_status": strings.NewReader("200"),
		"signature":             strings.NewReader(uploadPolicy.Signature),
		"callback":              strings.NewReader(uploadPolicy.Callback),
		"file":                  f.file,
	}
	return
}

// 核心函数 上传文件
// TODO: 上传进度条
func (f *UploadFile) upload() (err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range f.values {
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
			return err
		}

	}

	_ = w.Close()

	req, err := http.NewRequest("POST", f.url, &b)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	uploadBeginTime := time.Now()
	resp, err := f.client.Do(req)
	if err != nil {
		return
	}
	// 读取返回响应
	respBody, _ := ioutil.ReadAll(resp.Body)
	uploadUseTime := time.Since(uploadBeginTime)
	fmt.Println("upload use time:", uploadUseTime)

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("upload file to aliyun oss failed, status code:%s,response body: %s", resp.Status, string(respBody))
		return err
	}
	//  解析返回结果
	var uploadResponse UploadResponse
	err = json.Unmarshal(respBody, &uploadResponse)
	if err != nil {
		fmt.Println("handleUploadResponse json unmarshal failed", err)
		return err
	}
	f.uploadResponse = uploadResponse
	return nil
}
