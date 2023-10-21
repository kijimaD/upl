package cmd

import (
	"fmt"
	"io"
	"time"
)

var marks = []string{"|", "/", "-", "\\"}

func mark(i int) string {
	return marks[i%4]
}

func anim(w io.Writer) {
	go func() {
		i := 0
		for range time.Tick(100 * time.Millisecond) {
			if i == len(marks) {
				i = 0
			}
			fmt.Fprintf(w, "\rwait... %s", mark(i))
			i++
		}
	}()
}
