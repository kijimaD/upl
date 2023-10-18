package upl

import (
	"fmt"
	"io"
	"os/exec"
)

var (
	COMMAND       = "curl"
	BASEURL       = "http://localhost:7777"
	VERBOSE_OPT   = "--trace-ascii -"
	UPLOAD_TARGET = "upload.zip"
)

func buildUpload(cookie string) string {
	basecmd := `%s %s \
  %s \
  -# \
  -H 'Cookie: filemanager=%s' \
  --compressed \
  -F "p=" \
  -F "fullpath=%s" \
  -F "file=@%s;type=application/zip"`
	cmd := fmt.Sprintf(basecmd,
		COMMAND,
		BASEURL,
		VERBOSE_OPT,
		cookie,
		UPLOAD_TARGET,
		UPLOAD_TARGET,
	)
	return cmd
}

func Exec(out io.Writer) error {
	str, err := login(out)
	if err != nil {
		return err
	}
	cookie, err := parseCookie(str)
	if err != nil {
		return err
	}
	cmd := buildUpload(cookie)
	result, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		fmt.Fprint(out, string(result), err)
		return err
	}
	// TODO: 途中経過を表示したい
	return err
}
