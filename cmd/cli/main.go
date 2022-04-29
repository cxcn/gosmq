package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {

	if len(os.Args) == 1 {
		printInfo()
		fmt.Println("\nPress the enter key to exit...")
		fmt.Scanln()
	} else {
		cli()
	}
}

func printInfo() {
	fmt.Printf("gosmq-cli v0.23 %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("https://github.com/cxcn/gosmq/\n")
}
