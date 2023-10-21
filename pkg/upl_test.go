package upl

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

// クッキーが無効だとエラーを返す
func TestUploadFailInvalidCookie(t *testing.T) {
	file, err := os.Create(UPLOAD_TARGET)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())

	b := &bytes.Buffer{}
	task := NewTask(b)
	err = task.upload("invalid cookie")

	assert.Equal(t, FailUploadError, err)
}

// upload.zipファイルがないとエラーを返す
func TestUploadFailFileNotFound(t *testing.T) {
	b := &bytes.Buffer{}
	task := NewTask(b)
	cookie, err := task.login()
	if err != nil {
		t.Error(err)
	}
	err = task.upload(cookie)

	assert.Equal(t, TargetFileNotFoundError, err)
}
