package upload

import (
	"io/ioutil"
	"testing"
)

func Test_requestUploadPolicy(t *testing.T) {
	var u Instance
	resp, err := u.requestUploadPolicy()
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

func TestUploadInstance_buildUploadRequest(t *testing.T) {
	var u Instance
	err := u.getUploadPolicy()
	if err != nil {
		t.Error("buildUploadRequest error", err)
	}
	_, url, values, err := u.buildUploadRequest("/tmp/upload.txt")
	if err != nil {
		t.Error("buildUploadRequest error", err)
	}
	t.Logf("%s", url)
	t.Logf("%v", values)
}

func TestUploadInstance_handleUploadResponse(t *testing.T) {
	var u Instance
	uuidStr := "{\"uuid\":\"b5f55fdd-c71e-41d7-aead-5373e9196514\"}"
	u.handleUploadResponse([]byte(uuidStr))
}
