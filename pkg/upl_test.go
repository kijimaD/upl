package upl

import (
	"bytes"
	"os"
	"testing"
)

func TestExec(t *testing.T) {
	file, err := os.Create(UPLOAD_TARGET)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())

	b := &bytes.Buffer{}
	err = Exec(b)
	if err != nil {
		t.Error(err)
	}
}
