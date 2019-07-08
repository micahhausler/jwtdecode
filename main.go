package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/micahhausler/jwtdecode/pkg"
)

func main() {
	flag.Parse()
	filenames := flag.Args()

	if len(filenames) == 0 {
		filenames = []string{"/dev/stdin"}
	}
	err := pkg.DecodeFiles(os.Stdout, filenames)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
