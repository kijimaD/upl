package upl

import (
	"fmt"
	"io"
	"os/exec"
	"regexp"
)

func buildLogin() string {
	basecmd := `%s %s \
  -s \
  -c - \
  -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7' \
  -H 'Accept-Language: en-US,en;q=0.9,ja-JP;q=0.8,ja;q=0.7' \
  -H 'Cache-Control: max-age=0' \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -H 'Origin: http://localhost:7777' \
  -H 'Referer: http://localhost:7777/' \
  -H 'Sec-Fetch-Dest: document' \
  -H 'Sec-Fetch-Mode: navigate' \
  -H 'Sec-Fetch-Site: same-origin' \
  -H 'Sec-Fetch-User: ?1' \
  -H 'Upgrade-Insecure-Requests: 1' \
  -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36' \
  -H 'sec-ch-ua: "Google Chrome";v="117", "Not;A=Brand";v="8", "Chromium";v="117"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "Linux"' \
  --data-raw 'fm_usr=admin&fm_pwd=admin@123'`
	cmd := fmt.Sprintf(basecmd,
		COMMAND,
		BASEURL,
	)
	return cmd
}

// cookie jar textを返す
func login(out io.Writer) (string, error) {
	login, err := exec.Command("bash", "-c", buildLogin()).CombinedOutput()
	if err != nil {
		fmt.Fprint(out, login, err)
		return "", err
	}
	return string(login), nil
}

func parseCookie(in string) (string, error) {
	re, _ := regexp.Compile("\tfilemanager\t([0-9a-zA-Z]+)")
	ans := re.FindAllStringSubmatch(in, -1)
	return ans[0][1], nil
}
