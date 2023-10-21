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
