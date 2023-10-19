package upl

import (
	"fmt"
	"os/exec"
	"regexp"
)

const (
	ADMIN_USER = "admin"
	PWD        = "admin@123"
)

func (t *Task) buildLogin() string {
	basecmd := `%s %s \
  -s \
  -c - \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-raw 'fm_usr=%s&fm_pwd=%s'`
	cmd := fmt.Sprintf(basecmd,
		COMMAND,
		t.baseurl,
		ADMIN_USER,
		PWD,
	)
	return cmd
}

// cookie jar textを返す
func (t *Task) login() (string, error) {
	login, err := exec.Command("bash", "-c", t.buildLogin()).CombinedOutput()
	if err != nil {
		fmt.Fprint(t.w, string(login), err)
		return "", err
	}
	return string(login), nil
}

func (t *Task) parseCookie(in string) (string, error) {
	re, _ := regexp.Compile("\tfilemanager\t([0-9a-zA-Z]+)")
	ans := re.FindAllStringSubmatch(in, -1)
	return ans[0][1], nil
}
