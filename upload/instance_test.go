package upload

import (
	"io/ioutil"
	"testing"
)

func TestInstance_getCSRF(t *testing.T) {
	var u Instance
	err := u.getCSRF()
	if err != nil {
		t.Error("getCSRF error", err)
	}
	t.Log(u.CSRF)
	t.Log(u.Cookie)
}

func Test_requestUploadPolicy(t *testing.T) {
	resp, err := requestUploadPolicy()
	if err != nil {
		t.Error("requestUploadPolicy error", err)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	t.Log("response Status : ", resp.Status)
	t.Log("response Headers : ", resp.Header)
	t.Log("response Body : ", string(respBody))
}

func Test_getUploadPolicy(t *testing.T) {
	var u Instance
	err := u.getUploadPolicy()
	if err != nil {
		t.Error("requestUploadPolicy error", err)
	}
	t.Log("policy AccessID", u.UploadPolicy.AccessID)
	t.Log("policy Callback", u.UploadPolicy.Callback)
	t.Log("policy Dir", u.UploadPolicy.Dir)
	t.Log("policy Expire", u.UploadPolicy.Expire)
	t.Log("policy Host", u.UploadPolicy.Host)
	t.Log("policy Policy", u.UploadPolicy.Policy)
	t.Log("policy Signature", u.UploadPolicy.Signature)
}

func TestInstance_ready(t *testing.T) {
	var u Instance
	err := u.getUploadPolicy()
	if err != nil {
		t.Error("getUploadPolicy error", err)
	}
	err = u.ready([]string{"/tmp/test.txt"})
	if err != nil {
		t.Error("ready error", err)
	}
	for k, v := range u.UploadFiles[0].Values {
		t.Log(k, v)
	}
}

func TestInstance_run(t *testing.T) {
	var u Instance
	err := u.getUploadPolicy()
	if err != nil {
		t.Error("getUploadPolicy error", err)
	}
	err = u.ready([]string{"/tmp/test.txt"})
	if err != nil {
		t.Error("ready error", err)
	}
	errors := u.run()
	if len(errors) > 0 {
		for _, err := range errors {
			t.Error("run error", err)
		}
	}
	for _, file := range u.UploadFiles {
		t.Log(file.StatusCode)
		t.Log(file.Values)
		t.Log(file.UploadResponse)
	}
}

// 全流程测试
func TestInstance_finish(t *testing.T) {
	var u Instance
	err := u.getCSRF()
	if err != nil {
		t.Error("getCSRF error ", err)
	}
	err = u.getUploadPolicy()
	if err != nil {
		t.Error("getUploadPolicy error", err)
	}
	err = u.ready([]string{"/tmp/test.txt"})
	if err != nil {
		t.Error("ready error", err)
	}
	errors := u.run()
	if len(errors) > 0 {
		for _, err := range errors {
			t.Error("run error", err)
		}
	}
	err = u.finish()
	if err != nil {
		t.Error("finish error", err)
	}

	t.Log("Recode :", u.Recode)
	t.Log("Owner :", u.Owner)
}
