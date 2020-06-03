package upload

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/arcosx/Nuwa/util"
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

// 自定义修改
func (u *Instance) modify(newRecode string) []error {
	var errors []error
	if newRecode != "" {
		var changeRecode ChangeRecode
		changeRecode.Auth = u.Password
		changeRecode.NewReCode = newRecode
		changeRecode.UserToken = u.Owner
		err := requestChangeRecode(u.Recode, changeRecode, u.CSRF, u.Cookie)
		if err != nil {
			fmt.Println("modify your recode failed. use random recode instead")
			errors = append(errors, fmt.Errorf("modify user recode failed"))
		} else {
			// 成功修改值
			u.Recode = newRecode
		}
	}

	if u.Password != "" {
		var changePassword ChangePassword
		changePassword.UserToken = u.Owner
		changePassword.Auth = getSha256(u.Password)
		err := requestChangePassword(u.Recode, changePassword, u.CSRF, u.Cookie)
		if err != nil {
			fmt.Println("modify your password failed. use empty password instead")
			errors = append(errors, fmt.Errorf("modify your password failed"))
			u.Password = util.DEFAULT_PASSWORD //设定失败重置为默认值
		}
	} else {
		u.Password = util.DEFAULT_PASSWORD
	}
	if u.DownCount != 0 {
		var changeDownCount ChangeDownCount
		changeDownCount.UserToken = u.Owner
		changeDownCount.NewDownCount = u.DownCount
		err := requestDownCount(u.Recode, changeDownCount, u.CSRF, u.Cookie)
		if err != nil {
			fmt.Println("modify your down count failed. use default down count (10 times) instead")
			errors = append(errors, fmt.Errorf("modify your down count"))
			u.DownCount = util.DEFAULT_DOWN_COUNT
		}

	} else {
		u.DownCount = util.DEFAULT_DOWN_COUNT
	}
	if u.ExpireTime != 0 {
		var changeExpireTime ChangeExpireTime
		changeExpireTime.NewExpireTime = u.ExpireTime
		changeExpireTime.UserToken = u.Owner
		err := requestExpireTime(u.Recode, changeExpireTime, u.CSRF, u.Cookie)
		if err != nil {
			fmt.Println("modify your expire time failed. use default expire time (8 hour) instead")
			errors = append(errors, fmt.Errorf("modify your expire time"))
			u.ExpireTime = util.DEFAULT_EXPIRE_TIME
		}
	} else {
		u.ExpireTime = util.DEFAULT_EXPIRE_TIME
	}
	return errors
}

// 辅助函数

// 获取上传需要的策略参数
func requestUploadPolicy() (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", util.SAFEU_BASE_URL+"/v1/upload/policy", nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("requestUploadPolicy error", err)
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

func requestChangeRecode(oldRecode string, changeRecode ChangeRecode, csrfToken string, cookie string) (err error) {

	jsonStr, err := json.Marshal(changeRecode)
	if err != nil {
		fmt.Println("requestChangeRecode json marshal error", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", util.SAFEU_BASE_URL, "/v1/recode/", oldRecode), bytes.NewReader(jsonStr))
	if err != nil {
		fmt.Println("requestChangeRecode request create NewRequest failed", err)
		return err
	}
	req.Header.Set("x-csrf-token", csrfToken)
	req.Header.Set("cookie", cookie)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("origin", "https://safeu.a2os.club")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("requestChangeRecode call failed", err)
		return err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		// 解析错误码和错误信息
		var changeRequestErrorResponse ChangeRequestErrorResponse
		err = json.Unmarshal(respBody, &changeRequestErrorResponse)
		if err != nil {
			fmt.Println("requestChangeRecode changeRequestErrorResponse json unmarshal failed", err)
		}
		fmt.Println("requestChangeRecode response show some problem : ", changeRequestErrorResponse.Message)
		return fmt.Errorf("requestChangeRecode reponse error %s", respBody)

	}

	var changeRequestResponse ChangeRequestResponse
	err = json.Unmarshal(respBody, &changeRequestResponse)
	if err != nil {
		fmt.Println("requestChangeRecode changeRequestResponse json unmarshal failed", err)
		return err
	}
	if changeRequestResponse.Message != "ok" {
		fmt.Println("requestChangeRecode response show some problem : ", respBody)
		return fmt.Errorf("requestChangeRecode reponse error %s", respBody)
	}
	return nil
}

func requestChangePassword(recode string, changePassword ChangePassword, csrfToken string, cookie string) (err error) {

	jsonStr, err := json.Marshal(changePassword)
	if err != nil {
		fmt.Println("requestChangePassword json marshal error", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", util.SAFEU_BASE_URL, "/v1/password/", recode), bytes.NewReader(jsonStr))
	if err != nil {
		fmt.Println("requestChangePassword request create NewRequest failed", err)
		return err
	}
	req.Header.Set("x-csrf-token", csrfToken)
	req.Header.Set("cookie", cookie)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("origin", "https://safeu.a2os.club")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("requestChangePassword call failed", err)
		return err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		// 解析错误码和错误信息
		var changeRequestErrorResponse ChangeRequestErrorResponse
		err = json.Unmarshal(respBody, &changeRequestErrorResponse)
		if err != nil {
			fmt.Println("requestChangePassword changeRequestErrorResponse json unmarshal failed", err)
		}
		fmt.Println("requestChangePassword response show some problem : ", changeRequestErrorResponse.Message)
		return fmt.Errorf("requestChangePassword reponse error %s", respBody)

	}

	var changeRequestResponse ChangeRequestResponse
	err = json.Unmarshal(respBody, &changeRequestResponse)
	if err != nil {
		fmt.Println("requestChangePassword changeRequestResponse json unmarshal failed", err)
		return err
	}
	if changeRequestResponse.Message != "ok" {
		fmt.Println("requestChangePassword response show some problem : ", respBody)
		return fmt.Errorf("requestChangePassword reponse error %s", respBody)
	}
	return nil
}

func requestExpireTime(recode string, changeExpireTime ChangeExpireTime, csrfToken string, cookie string) (err error) {

	jsonStr, err := json.Marshal(changeExpireTime)
	if err != nil {
		fmt.Println("requestExpireTime json marshal error", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", util.SAFEU_BASE_URL, "/v1/expireTime/", recode), bytes.NewReader(jsonStr))
	if err != nil {
		fmt.Println("requestExpireTime request create NewRequest failed", err)
		return err
	}
	req.Header.Set("x-csrf-token", csrfToken)
	req.Header.Set("cookie", cookie)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("origin", "https://safeu.a2os.club")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("requestExpireTime call failed", err)
		return err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		// 解析错误码和错误信息
		var changeRequestErrorResponse ChangeRequestErrorResponse
		err = json.Unmarshal(respBody, &changeRequestErrorResponse)
		if err != nil {
			fmt.Println("requestExpireTime changeRequestErrorResponse json unmarshal failed", err)
		}
		fmt.Println("requestExpireTime response show some problem : ", changeRequestErrorResponse.Message)
		return fmt.Errorf("requestExpireTime reponse error %s", string(respBody))

	}

	var changeRequestResponse ChangeRequestResponse
	err = json.Unmarshal(respBody, &changeRequestResponse)
	if err != nil {
		fmt.Println("requestExpireTime changeRequestResponse json unmarshal failed", err)
		return err
	}
	if changeRequestResponse.Message != "ok" {
		fmt.Println("requestExpireTime response show some problem : ", respBody)
		return fmt.Errorf("requestExpireTime reponse error %s", respBody)
	}
	return nil
}

func requestDownCount(recode string, changeDownCount ChangeDownCount, csrfToken string, cookie string) (err error) {

	jsonStr, err := json.Marshal(changeDownCount)
	if err != nil {
		fmt.Println("requestDownCount json marshal error", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", util.SAFEU_BASE_URL, "/v1/downCount/", recode), bytes.NewReader(jsonStr))
	if err != nil {
		fmt.Println("requestDownCount request create NewRequest failed", err)
		return err
	}
	req.Header.Set("x-csrf-token", csrfToken)
	req.Header.Set("cookie", cookie)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("origin", "https://safeu.a2os.club")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("requestDownCount call failed", err)
		return err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		// 解析错误码和错误信息
		var changeRequestErrorResponse ChangeRequestErrorResponse
		err = json.Unmarshal(respBody, &changeRequestErrorResponse)
		if err != nil {
			fmt.Println("requestDownCount changeRequestErrorResponse json unmarshal failed", err)
		}
		fmt.Println("requestDownCount response show some problem : ", changeRequestErrorResponse.Message)
		return fmt.Errorf("requestDownCount reponse error %s", string(respBody))

	}

	var changeRequestResponse ChangeRequestResponse
	err = json.Unmarshal(respBody, &changeRequestResponse)
	if err != nil {
		fmt.Println("requestDownCount changeRequestResponse json unmarshal failed", err)
		return err
	}
	if changeRequestResponse.Message != "ok" {
		fmt.Println("requestDownCount response show some problem : ", respBody)
		return fmt.Errorf("requestDownCount reponse error %s", string(respBody))
	}
	return nil
}

func getSha256(text string) string {
	bv := []byte(text)
	hasher := sha256.New()
	hasher.Write(bv)
	sha := hex.EncodeToString(hasher.Sum(nil))
	return sha
}

// 入口函数
func Start(fileFullPaths []string, userRecode string, userPassword string, userDownCount int, userExpireTime int) {
	var u Instance
	u.Password = userPassword
	u.DownCount = userDownCount
	u.ExpireTime = userExpireTime

	err := u.getCSRF()
	if err != nil {
		fmt.Println("getCSRF error", err)
		os.Exit(0)
	}
	err = u.getUploadPolicy()
	if err != nil {
		fmt.Println("getUploadPolicy error", err)
		os.Exit(0)
	}
	err = u.ready(fileFullPaths)
	if err != nil {
		fmt.Println("ready error", err)
		os.Exit(0)
	}
	errors := u.run()
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println("file upload error", err)
		}
		os.Exit(0)
	}
	err = u.finish()
	if err != nil {
		fmt.Println("finish error", err)
		os.Exit(0)
	}
	errors = u.modify(userRecode)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println("modify error", err)
		}
	}
	fmt.Println("Upload Finish")
	fmt.Println("")

	fmt.Println("Recode :", u.Recode)
	fmt.Println("Owner :", u.Owner)
	fmt.Println("Password :", u.Password)
	fmt.Println("DownCount :", u.DownCount)
	fmt.Println("ExpireTime :", u.ExpireTime)
}
