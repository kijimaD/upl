package main

import (
	"fmt"
	"os/exec"
)

// TODO: クッキーを取ってくる

var cmd = `
curl 'http://localhost:7777?p=' \
  -v \
  --trace-ascii - \
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
  -F "file=@upload.zip;type=application/zip"
`

func main() {
	result, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(result))
}
