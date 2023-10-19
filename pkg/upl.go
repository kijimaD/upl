package upl

import (
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

	go func() {
		t.mu.Lock()
		io.Copy(t.w, stderr) // curlのプログレスバーはなぜか標準エラー出力である
		t.mu.Unlock()
		// curl の --limit-rate 1m オプションで転送速度を遅くして動作確認できる
	}()

	err = cmd.Wait()
	if err != nil {
		return err
	}
	t.mu.Lock()
	io.Copy(t.w, stdout)
	t.mu.Unlock()
	return err
}
