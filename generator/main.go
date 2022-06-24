package main

import (
	"github.com/apache/incubator-devlake/generator/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
