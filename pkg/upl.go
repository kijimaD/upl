package upl

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"sync"
)

var (
	DEFAULT_BASEURL = "http://localhost:7777"
	COMMAND         = "curl"
	VERBOSE_OPT     = "--trace-ascii -"
	UPLOAD_TARGET   = "upload.zip"

	TargetFileNotFoundError = errors.New("カレントディレクトリに upload.zip ファイルが存在しない")
	FailLoginError          = errors.New("ユーザー認証に失敗した")
	FailUploadError         = errors.New("アップロードが成功レスポンスを返却しなかった")
)

const (
	DEFAULT_ADMINUSER = "admin"
	DEFAULT_PWD       = "admin@123"
)

type Task struct {
	w         io.Writer
	mu        *sync.RWMutex
	baseurl   string
	adminuser string
	pwd       string
	destpath  string
}

func NewTask(w io.Writer, options ...TaskOption) *Task {
	task := Task{
		w:         w,
		mu:        &sync.RWMutex{},
		baseurl:   DEFAULT_BASEURL,
		adminuser: DEFAULT_ADMINUSER,
		pwd:       DEFAULT_PWD,
		destpath:  ".",
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

func TaskWithLoginUser(adminuser string, pwd string) TaskOption {
	return func(t *Task) {
		t.adminuser = adminuser
		t.pwd = pwd
	}
}

func TaskWithDestpath(destpath string) TaskOption {
	return func(t *Task) {
		t.destpath = destpath
	}
}

func (t *Task) upload(cookie string) error {
	file, err := os.Open(UPLOAD_TARGET)
	if err != nil {
		return TargetFileNotFoundError
	}
	defer file.Close()

	boundary := "---------------------------1234567890"

	// req body作成
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.SetBoundary(boundary)
	writer.WriteField("p", t.destpath)
	writer.WriteField("fullpath", UPLOAD_TARGET)
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, UPLOAD_TARGET))
	partHeader.Set("Content-Type", "application/zip")
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	writer.Close()

	// リクエストを作成
	req, err := http.NewRequest("POST", t.baseurl, body)
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", fmt.Sprintf("filemanager=%s", cookie))
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = int64(body.Len())

	// リクエストを送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(bodyText) != `{"status":"success","info":"file upload successful"}` {
		return FailUploadError
	}
	return nil
}

func (t *Task) Exec() error {
	_, err := os.Stat(UPLOAD_TARGET)
	if err != nil {
		return TargetFileNotFoundError
	}
	cookie, err := t.login()
	if err != nil {
		return err
	}
	err = t.upload(cookie)
	if err != nil {
		return err
	}
	return err
}

// curl http://localhost:7777
// -#
// -H
// "Cookie: filemanager=gbpenh0c3r0ucvr2ijmmflact6"
// --compressed
// -F "p="
// -F "fullpath=upload.zip"
// -F "file=@upload.zip;type=application/zip"

// 0000: POST / HTTP/1.1
// 0011: Host: localhost:7777
// 0027: User-Agent: curl/7.85.0
// 0040: Accept: */*
// 004d: Accept-Encoding: deflate, gzip
// 006d: Cookie: filemanager=ihgg0r646rqe1iei38bgjksp0s
// 009d: Content-Length: 390
// 00b2: Content-Type: multipart/form-data; boundary=--------------------
// 00f2: ----b290fafbc9171b05
// 0108:
// => Send data, 390 bytes (0x186)
// 0000: --------------------------b290fafbc9171b05
// 002c: Content-Disposition: form-data; name="p"
// 0056:
// 0058:
// 005a: --------------------------b290fafbc9171b05
// 0086: Content-Disposition: form-data; name="fullpath"
// 00b7:
// 00b9: upload.zip
// 00c5: --------------------------b290fafbc9171b05
// 00f1: Content-Disposition: form-data; name="file"; filename="upload.zi
// 0131: p"
// 0135: Content-Type: application/zip
// 0154:
// 0156:
// 0158: --------------------------b290fafbc9171b05--
