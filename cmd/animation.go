package cmd

import (
	"fmt"
	"io"
	"time"
)

var marks = []string{"|", "/", "-", "\\"}

const animLoopMilliseconds = 100

func mark(i int) string {
	return marks[i%len(marks)]
}

func anim(w io.Writer) {
	go func() {
		i := 0
		for range time.Tick(animLoopMilliseconds * time.Millisecond) {
			if i == len(marks) {
				i = 0
			}
			fmt.Fprintf(w, "\rwait... %s", mark(i))
			i++
		}
	}()
}
