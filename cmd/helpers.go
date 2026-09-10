package cmd

import "github.com/spf13/cobra"

func walkCommandTree(root *cobra.Command, fn func(cmd *cobra.Command)) {
	fn(root)
	for _, child := range root.Commands() {
		walkCommandTree(child, fn)
	}
}

func walkCommandTreeWithFlag(root *cobra.Command, flag string, fn func(cmd *cobra.Command)) {
	walkCommandTree(root, func(cmd *cobra.Command) {
		if cmd.Flags().Lookup(flag) == nil {
			return
		}
		fn(cmd)
	})
}
