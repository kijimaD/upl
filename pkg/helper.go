package upl

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

func displayOutput(r io.Reader, w io.Writer) {
	const timerDisplaySyncSec = 0
	scanner := bufio.NewScanner(r)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				scannedText := scanner.Text()

				if len(scannedText) > 0 {
					head := fmt.Sprintf("%s", scannedText)
					fmt.Fprintf(w, "\r%s", head)
				}
				time.Sleep(timerDisplaySyncSec * time.Millisecond)
			}
		}
	}()

	// 行を次に進める
	for scanner.Scan() {
		scannedText := scanner.Text()
		head := fmt.Sprintf("%s", scannedText)
		fmt.Fprintf(w, "%s\n", head)
	}
	done <- true
}
