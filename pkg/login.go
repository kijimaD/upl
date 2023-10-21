package upl

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

const (
	ADMIN_USER = "admin"
	PWD        = "admin@123"
)

func (t *Task) buildLogin() string {
	basecmds := []string{
		"%s",
		"%s",
		"-s",
		"-c -",
		"-H 'Content-Type: application/x-www-form-urlencoded'",
		"--data-raw 'fm_usr=%s&fm_pwd=%s'",
	}
	basecmd := strings.Join(basecmds, " ")
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
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "-Command", t.buildLogin())
	default: // Linux & Mac
		cmd = exec.Command("sh", "-c", t.buildLogin())
	}
	login, err := cmd.CombinedOutput()
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

// ログインが成功しない
func (t *Task) login2() (string, error) {
	form := url.Values{}
	form.Add("fm_usr", "admin")
	form.Add("fm_pwd", "admin@123")
	body := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", t.baseurl, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Go/1.20")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	// クッキーの処理: クッキーを格納するためのジャーを作成
	// jar, _ := cookiejar.New(nil)
	// client.Jar = jar

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	var c string
	for _, cookie := range cookies {
		if cookie.Name == "filemanager" {
			c = cookie.Value
		}
	}
	// ログインに失敗するとトップページを返す
	if len(c) == 0 {
		return "", errors.New("クッキーを取得できなかった")
	}
	if resp.Status != "302 Found" {
		return "", errors.New("ログインが成功しなかった")
	}

	dump, _ := httputil.DumpRequest(req, true)
	fmt.Printf("%s", dump)
	dump2, _ := httputil.DumpResponse(resp, true)
	fmt.Printf("%s", dump2)

	// ログインが成功してない
	return c, nil
}
