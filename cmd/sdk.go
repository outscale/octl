package cmd

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/smithy-go/logging"
	"github.com/outscale/octl/pkg/messages"
	"github.com/outscale/octl/pkg/sdk"
	"github.com/outscale/octl/pkg/version"
	"github.com/outscale/osc-sdk-go/v3/pkg/middleware"
	"github.com/outscale/osc-sdk-go/v3/pkg/oos"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	"github.com/spf13/cobra"
)

func loadProfile(cmd *cobra.Command) profile.Result {
	var opt profile.Options

	path, _ := cmd.Flags().GetString("config")
	if path != "" {
		opt.FilePath = &path
	}

	prof, _ := cmd.Flags().GetString("profile")
	if prof != "" {
		opt.ProfileName = &prof
	}

	p, err := profile.Load(opt)
	if err != nil {
		messages.ExitErr(err)
	}

	return p
}

func sdkOptions(cmd *cobra.Command) []middleware.MiddlewareChainOption {
	ua := "octl/" + version.Version
	opts := []middleware.MiddlewareChainOption{options.WithUseragent(ua)}
	if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
		opts = append(opts, options.WithLogging(sdk.VerboseLogger{}))
	} else {
		opts = append(opts, options.WithoutLogging())
	}
	return opts
}

func awsOptions(cmd *cobra.Command) []config.LoadOptionsFunc {
	ua := "octl/" + version.Version
	opts := []config.LoadOptionsFunc{config.WithAppID(ua)}
	if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
		opts = append(opts,
			config.WithClientLogMode(aws.LogRequest|aws.LogRequestWithBody|aws.LogResponseWithBody),
			config.WithLogger(logging.NewStandardLogger(os.Stderr)),
			oos.WithUseragent(ua),
		)
	}
	return opts
}
