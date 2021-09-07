package main

import (
	"log"

	"github.com/merico-dev/lake/cmd/api"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:     "lake-cli",
	Short:   "Lake cli tools",
	Version: "v0.0.1",
}

func main() {
	api.Register(root)
	if err := root.Execute(); err != nil {
		log.Fatalln(err)
	}
}
