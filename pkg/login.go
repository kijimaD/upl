package upl

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	ADMIN_USER = "admin"
	PWD        = "admin@123"
)

// クッキーを生成する
func (t *Task) getCookie() (string, error) {
	req, err := http.NewRequest("GET", t.baseurl, nil)
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
func (t *Task) login() (string, error) {
	values := url.Values{}
	values.Add("fm_usr", ADMIN_USER)
	values.Add("fm_pwd", PWD)

	req, err := http.NewRequest("POST", t.baseurl, strings.NewReader(values.Encode()))
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
	if resp.StatusCode != http.StatusOK {
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
