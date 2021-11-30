package main

import (
	"os"
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
