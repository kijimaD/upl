package upl

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ログインできる
func TestLogin(t *testing.T) {
	buf := &bytes.Buffer{}
	task := NewTask(buf)
	str, err := task.login()
	if err != nil {
		t.Error(err)
	}
	cookie, err := task.parseCookie(str)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 26, len(cookie))
}

func TestLogin2(t *testing.T) {
	buf := &bytes.Buffer{}
	task := NewTask(buf)
	cookie, err := task.login2()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 26, len(cookie))
}
