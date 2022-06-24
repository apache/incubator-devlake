package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func init() {
	rootCmd.AddCommand(GeneratorDocCmd)
}

var GeneratorDocCmd = &cobra.Command{
	Use:   "generator-doc",
	Short: "generate document for generator",
	RunE: func(cmd *cobra.Command, args []string) error {
		return doc.GenMarkdownTree(rootCmd, "generator/docs")
	},
}
