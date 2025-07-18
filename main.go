package main

import (
	"fmt"
	"go-api-gen/cmd"
	"os"
)

var Version = "v1.0.3"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println(Version)
		return
	}
	cmd.Execute()
}
