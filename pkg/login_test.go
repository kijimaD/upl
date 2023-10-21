package upl

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	buf := &bytes.Buffer{}
	task := NewTask(buf)
	cookie, err := task.login()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 26, len(cookie))
}

func TestLoginFail(t *testing.T) {
	buf := &bytes.Buffer{}
	task := NewTask(buf, TaskWithLoginUser("invalid user", "invalid pwd"))
	cookie, err := task.login()
	// 認証情報が無効だとエラーを返す
	if err != FailLoginError {
		t.Error(err)
	}
	assert.Equal(t, 0, len(cookie))
}
