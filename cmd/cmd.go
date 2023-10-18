package cmd

import (
	"flag"
	"io"

	upl "github.com/kijimaD/upl/pkg"
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
	err := upl.Exec()
	if err != nil {
		return err
	}

	return nil
}
