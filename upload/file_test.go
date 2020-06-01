package upload

import (
	"testing"
)

func TestUploadFile_buildUploadRequest(t *testing.T) {
	var u Instance
	var f UploadFile
	err := u.getUploadPolicy()
	if err != nil {
		t.Error("buildUploadRequest error", err)
	}
	err = f.buildUploadRequest(u.UploadPolicy, "/tmp/upload.txt")
	if err != nil {
		t.Error("buildUploadRequest error", err)
	}
	t.Logf("%s", f.url)
	t.Logf("%v", f.values)
}
