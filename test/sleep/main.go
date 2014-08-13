// sleep is a test program for testing Process
//
// sleep doesn't exist on Windows this is a substitute
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "secondCount")
		os.Exit(1)
	}

	n, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	time.Sleep(time.Duration(n) * time.Second)
}
