package upload

import (
	"io/ioutil"
	"testing"
)

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
