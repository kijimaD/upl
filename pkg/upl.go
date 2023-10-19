package upl

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

var (
	DEFAULT_BASEURL = "http://localhost:7777"
	COMMAND         = "curl"
	VERBOSE_OPT     = "--trace-ascii -"
	UPLOAD_TARGET   = "upload.zip"
)

type Task struct {
	w       io.Writer
	mu      *sync.RWMutex
	baseurl string
}

func NewTask(w io.Writer) *Task {
	task := Task{
		w:       w,
		mu:      &sync.RWMutex{},
		baseurl: DEFAULT_BASEURL,
	}
	return &task
}

func (t *Task) buildUpload(cookie string) string {
	basecmd := `%s %s \
  -# \
  -H 'Cookie: filemanager=%s' \
  --compressed \
  -F "p=" \
  -F "fullpath=%s" \
  -F "file=@%s;type=application/zip"`
	cmd := fmt.Sprintf(basecmd,
		COMMAND,
		t.baseurl,
		cookie,
		UPLOAD_TARGET,
		UPLOAD_TARGET,
	)
	return cmd
}

func (t *Task) Exec() error {
	str, err := t.login()
	if err != nil {
		return err
	}
	cookie, err := t.parseCookie(str)
	if err != nil {
		return err
	}
	cmdtext := t.buildUpload(cookie)
	cmd := exec.Command("bash", "-c", cmdtext)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	go t.displayOutput(stdout)
	go t.displayOutput(stderr)

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return err
}

func (t *Task) displayOutput(r io.Reader) {
	const timerDisplaySyncSec = 100
	scanner := bufio.NewScanner(r)
	done := make(chan bool)

	for scanner.Scan() {
		scannedText := scanner.Text()
		head := fmt.Sprintf("%s", scannedText)
		t.mu.Lock()
		fmt.Fprintf(t.w, "%s\n", head)
		t.mu.Unlock()
	}
	done <- true
}
