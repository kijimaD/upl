package upl

import (
	"fmt"
	"io"
	"os/exec"
	"regexp"
)

var (
	COMMAND     = "curl"
	BASEURL     = "http://localhost:7777"
	VERBOSE_OPT = "--trace-ascii -"
)

func build() string {
	basecmd := `%s %s \
  %s \
  -H 'Accept: application/json' \
  -H 'Accept-Language: ja,en-US;q=0.9,en;q=0.8' \
  -H 'Cache-Control: no-cache' \
  -H 'Connection: keep-alive' \
  -H 'Cookie: filemanager=fbd4ccd74ea741e10c489344382c8958' \
  -H 'Origin: http://localhost' \
  -H 'Pragma: no-cache' \
  -H 'Referer: http://localhost/uploader?p=&upload' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-origin' \
  -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36' \
  -H 'X-Requested-With: XMLHttpRequest' \
  -H 'sec-ch-ua: "Google Chrome";v="117", "Not;A=Brand";v="8", "Chromium";v="117"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "Linux"' \
  --compressed \
  -F "p=" \
  -F "fullpath=upload.zip" \
  -F "file=@upload.zip;type=application/zip"`
	cmd := fmt.Sprintf(basecmd,
		COMMAND,
		BASEURL,
		VERBOSE_OPT,
	)
	return cmd
}

func buildLogin() string {
	basecmd := `curl 'http://localhost:7777/' \
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
  --data-raw 'fm_usr=admin&fm_pwd=admin%40123'`
	return basecmd
}

func login(out io.Writer) error {
	login, err := exec.Command("bash", "-c", buildLogin()).CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Fprintln(out, string(login))
	return nil
}

func parseCookie(in string) (string, error) {
	re, _ := regexp.Compile("\tfilemanager\t([0-9a-zA-Z]+)$")
	ans := re.FindAllStringSubmatch(in, -1)
	return ans[0][1], nil
}

func Exec(out io.Writer) error {
	result, err := exec.Command("bash", "-c", build()).CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Fprintf(out, string(result))
	return err
}
