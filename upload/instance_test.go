package upload

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// use for generate a random string
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

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

func Test_getSha256(t *testing.T) {
	test := "hello SafeU"
	fmt.Println("78b12bba56d6f4a6b94faa89163994a14a92f2d246460751b4a48747fd90cf81", getSha256(test))
}

func TestStart(t *testing.T) {
	// create file for this test
	// please this action will rewrite your local file !
	err := ioutil.WriteFile("/tmp/testsafeucli.txt", []byte("Hello"), 0755)
	if err != nil {
		t.Error("Unable to write test file ", err)
	}
	err = ioutil.WriteFile("/tmp/testsafeucli2.txt", []byte("SafeU"), 0755)
	if err != nil {
		t.Error("Unable to write test file ", err)
	}
	type args struct {
		fileFullPaths  []string
		userRecode     string
		userPassword   string
		userDownCount  int
		userExpireTime int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "one file with out modify",
			args: args{
				fileFullPaths:  []string{"/tmp/testsafeucli.txt"},
				userRecode:     "",
				userPassword:   "",
				userDownCount:  0,
				userExpireTime: 0,
			},
		},
		{
			name: "more file with out modify",
			args: args{
				fileFullPaths:  []string{"/tmp/testsafeucli.txt", "/tmp/testsafeucli2.txt"},
				userRecode:     "",
				userPassword:   "",
				userDownCount:  0,
				userExpireTime: 0,
			},
		},
		{
			name: "one file with modify",
			args: args{
				fileFullPaths:  []string{"/tmp/testsafeucli.txt"},
				userRecode:     RandStringRunes(10),
				userPassword:   RandStringRunes(10),
				userDownCount:  4,
				userExpireTime: 10,
			},
		},
		{
			name: "more file with modify",
			args: args{
				fileFullPaths:  []string{"/tmp/testsafeucli.txt", "/tmp/testsafeucli2.txt"},
				userRecode:     RandStringRunes(10),
				userPassword:   RandStringRunes(10),
				userDownCount:  4,
				userExpireTime: 8,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Start(tt.args.fileFullPaths, tt.args.userRecode, tt.args.userPassword, tt.args.userDownCount, tt.args.userExpireTime)
		})
	}
}
