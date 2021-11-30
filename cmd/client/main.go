package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	// defer profile.Start().Stop()
	// defer profile.Start(profile.MemProfile, profile.MemProfileRate(1)).Stop()

	if len(os.Args) == 1 {
		web()
	} else {
		cli()
	}

	// time.Sleep(5 * time.Second)
}

func printInfo() {
	fmt.Printf("smq-client version 0.14 %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("repo address: https://github.com/cxcn/gosmq/\n\n")
}
