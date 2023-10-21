package upl

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var (
	DEFAULT_BASEURL = "http://localhost:7777"
	COMMAND         = "curl"
	VERBOSE_OPT     = "--trace-ascii -"
	UPLOAD_TARGET   = "upload.zip"

	TargetFileNotFoundError = errors.New("カレントディレクトリに upload.zip ファイルが存在しない")
)

type Task struct {
	w       io.Writer
	mu      *sync.RWMutex
	baseurl string
}

func NewTask(w io.Writer, options ...TaskOption) *Task {
	task := Task{
		w:       w,
		mu:      &sync.RWMutex{},
		baseurl: DEFAULT_BASEURL,
	}
	for _, option := range options {
		option(&task)
	}
	return &task
}

type TaskOption func(*Task)

func TaskWithBaseurl(baseurl string) TaskOption {
	return func(t *Task) {
		t.baseurl = baseurl
	}
}

func (t *Task) buildUpload(cookie string) string {
	basecmds := []string{
		`%s`,
		`%s`,
		`-#`,
		`-H "Cookie: filemanager=%s"`,
		`--compressed`,
		`-F "p="`,
		`-F "fullpath=%s"`,
		`-F "file=@%s;type=application/zip"`,
	}
	basecmd := strings.Join(basecmds, " ")
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
	_, err := os.Stat(UPLOAD_TARGET)
	if err != nil {
		return TargetFileNotFoundError
	}

	str, err := t.login()
	if err != nil {
		return err
	}
	cookie, err := t.parseCookie(str)
	if err != nil {
		return err
	}
	cmdtext := t.buildUpload(cookie)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "-Command", cmdtext)
	default: // Linux & Mac
		cmd = exec.Command("sh", "-c", cmdtext)
	}

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
