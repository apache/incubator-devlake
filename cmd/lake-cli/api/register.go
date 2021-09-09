package api

import "github.com/spf13/cobra"

// Register registers all commands.
func Register(root *cobra.Command) {
	root.AddCommand(newApiCommand())
}
