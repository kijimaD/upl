package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"

	upl "github.com/kijimaD/upl/pkg"
)

var (
	ArgumentCountError = errors.New("引数の数が間違っている。expect: 1")
)

type CLI struct {
	Out io.Writer
}

func New(out io.Writer) *CLI {
	return &CLI{
		Out: out,
	}
}

func (cli *CLI) Execute(args []string) error {
	flag.Parse()

	// 何も指定されてないときはヘルプを出す
	help := `Usage:
upl {Tiny File Manager base URL}

uplはカレントディレクトリにある "upload.zip" という名前のファイルを、
指定したパスにあるTiny File Managerに転送します。
`
	if len(args) == 1 {
		fmt.Fprintf(cli.Out, "%s\n", help)
		return nil
	}

	if len(args) != 2 {
		return ArgumentCountError
	}
	baseurl := args[1]

	anim(cli.Out)

	task := upl.NewTask(cli.Out, upl.TaskWithBaseurl(baseurl))
	err := task.Exec()
	if err != nil {
		return err
	}

	fmt.Fprintln(cli.Out, "")
	return nil
}
