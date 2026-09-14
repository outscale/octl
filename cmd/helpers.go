package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

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

func setFlag(f *pflag.Flag, v string) error {
	f.Changed = true
	return f.Value.Set(v)
}
