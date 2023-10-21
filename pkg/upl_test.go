package upl

import (
	"bytes"
	"os"
	"testing"
)

func TestUpload(t *testing.T) {
	file, err := os.Create(UPLOAD_TARGET)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())

	b := &bytes.Buffer{}
	task := NewTask(b)
	cookie, err := task.login()
	if err != nil {
		t.Error(err)
	}
	err = task.upload(cookie)
	if err != nil {
		t.Error(err)
	}
}

func TestUploadFail(t *testing.T) {
	file, err := os.Create(UPLOAD_TARGET)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())

	b := &bytes.Buffer{}
	task := NewTask(b)
	err = task.upload("invalid cookie")
	// クッキーが無効だと失敗する
	if err != FailUploadError {
		t.Error(err)
	}
}
