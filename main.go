package main

import (
	"os"

	"github.com/kijimaD/upl/cmd"
)

func main() {
	cmd := cmd.New(os.Stdout)
	err := cmd.Execute(os.Args)
	if err != nil {
		panic(err)
	}
}
