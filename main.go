package main

import (
	"flag"
	"fmt"
	"go-api-gen/cmd"
)

func main() {
	flagVersion := flag.Bool("version", false, "Mostra a versão")
	flag.Parse()

	if *flagVersion {
		fmt.Println("go-api-gen v1.0.5")
		return
	}
	cmd.Execute()
}
