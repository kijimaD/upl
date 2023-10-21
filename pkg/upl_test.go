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
	resp, err := task.upload(cookie)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, `{"status":"success","info":"file upload successful"}`, resp)
}
