package commandbuilder

import (
	"github.com/outscale/octl/pkg/debug"
	"github.com/spf13/cobra"
)

// BuildAPI builds the low-level API.
func (b *Builder) BuildAPI(rootCmd *cobra.Command, run func(cmd *cobra.Command, args []string)) *cobra.Command {
	rootCmd.AddGroup(&cobra.Group{
		ID:    "api",
		Title: "API",
	})
	apiCmd := &cobra.Command{
		Use:     "api",
		GroupID: "api",
		Short:   "Call " + rootCmd.Use + " API",
	}
	rootCmd.AddCommand(apiCmd)

	for _, call := range b.cfg.API {
		if call.Use == "" {
			continue
		}
		if call.Group != "" && !apiCmd.ContainsGroup(call.Group) {
			apiCmd.AddGroup(&cobra.Group{ID: call.Group, Title: call.Group})
		}
		cmd := &cobra.Command{
			// we need to disable flag parsing, an alias might generate flags
			// that were unknown during init
			DisableFlagParsing: true,
			GroupID:            call.Group,
			Use:                call.Use,
			Short:              call.Short,
			Long:               call.Help,
			Run:                run,
			PreRun: func(cmd *cobra.Command, args []string) {
				if !cmd.DisableFlagParsing {
					return
				}
				// update flag set with args
				// a resolved alias will generate a new command line with potential new flags.
				debug.Println("Updating flags for", cmd.Name())
				err := b.buildFlagSet(cmd, call.Flags, getNumEntriesInSlices(args))
				if err != nil {
					debug.Println(call.Entity, call.Use, err)
				}
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				if cmd.DisableFlagParsing {
					// loop with enabled flag parsing
					debug.Println("Looping on", cmd.Use, "with", args)
					cmd.DisableFlagParsing = false
					return cmd.Execute()
				}
				debug.Println("Running", cmd.Use, "with", args)
				run(cmd, args)
				return nil
			},
		}
		apiCmd.AddCommand(cmd)
		// build default flag set based on CLI args
		err := b.buildFlagSet(cmd, call.Flags, osNumEntriesInSlices)
		if err != nil {
			debug.Println(call.Entity, call.Use, err)
		}
	}
	return apiCmd
}
