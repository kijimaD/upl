package upl

import (
	"errors"
	"fmt"
	"net/http"
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

// クッキーを生成する
func (t *Task) getCookie() (string, error) {
	req, err := http.NewRequest("POST", "http://localhost:7777/index.php?p=", nil)
	client := &http.Client{}

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
	return c, nil
}

// 生成したクッキーでログインする
func (t *Task) login2() (string, error) {
	values := url.Values{}
	values.Add("fm_usr", "admin")
	values.Add("fm_pwd", "admin@123")

	req, err := http.NewRequest("POST", "http://localhost:7777/index.php?p=", strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", "Go/1.20")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	cookie, err := t.getCookie()
	if err != nil {
		return "", err
	}
	if len(cookie) == 0 {
		return "", errors.New("クッキーを取得できなかった")
	}
	req.Header.Set("Cookie", fmt.Sprintf("filemanager=%s", cookie))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return "", errors.New("ログインが成功しなかった")
	}

	// dump1, _ := httputil.DumpRequest(req, true)
	// fmt.Printf("%s", dump1)
	// dump2, _ := httputil.DumpResponse(resp, true)
	// fmt.Printf("%s", dump2)

	return cookie, nil
}

// curl 'http://localhost:7777/index.php?p=' \
//   -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7' \
//   -H 'Accept-Language: en-US,en;q=0.9,ja-JP;q=0.8,ja;q=0.7' \
//   -H 'Cache-Control: max-age=0' \
//   -H 'Connection: keep-alive' \
//   -H 'Content-Type: application/x-www-form-urlencoded' \
//   -H 'Cookie: filemanager=3ovivhi019m5313hdnmolqj7ku' \
//   -H 'Origin: http://localhost:7777' \
//   -H 'Referer: http://localhost:7777/index.php?p=' \
//   -H 'Sec-Fetch-Dest: document' \
//   -H 'Sec-Fetch-Mode: navigate' \
//   -H 'Sec-Fetch-Site: same-origin' \
//   -H 'Sec-Fetch-User: ?1' \
//   -H 'Upgrade-Insecure-Requests: 1' \
//   -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36' \
//   -H 'sec-ch-ua: "Google Chrome";v="117", "Not;A=Brand";v="8", "Chromium";v="117"' \
//   -H 'sec-ch-ua-mobile: ?0' \
//   -H 'sec-ch-ua-platform: "Linux"' \
//   --data-raw 'fm_usr=admin&fm_pwd=admin%40123' \
//   --compressed
