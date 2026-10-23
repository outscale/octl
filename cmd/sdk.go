package cmd

import (
	"errors"
	"os"
	"slices"
	"strings"

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
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type namedProfile struct {
	name string
	p    *profile.Profile
}

// loadProfiles resolves the --profile flag to one or more profiles:
// "all" or a comma-separated list of names selects several profiles
// from the profile file, any other value keeps the single profile behavior.
func loadProfiles(cmd *cobra.Command) ([]namedProfile, error) {
	value, _ := cmd.Flags().GetString("profile")
	names := strings.Split(value, ",")
	if len(names) == 1 && names[0] != "all" {
		return []namedProfile{{name: value, p: loadProfile(cmd)}}, nil
	}
	cf, err := loadConfig(cmd)
	if err != nil {
		return nil, err
	}
	if names[0] == "all" {
		names = lo.Keys(cf.Profiles)
		slices.Sort(names)
	}
	profs := make([]namedProfile, 0, len(names))
	for _, name := range names {
		if name == "" {
			continue
		}
		raw, err := cf.Profile(name)
		if err != nil {
			return nil, err
		}
		// inject the profile as an option to get the same defaults as profile.New
		p, err := profile.New(func(pp *profile.Profile) error { *pp = raw; return nil })
		if err != nil {
			return nil, err
		}
		profs = append(profs, namedProfile{name: name, p: p})
	}
	if len(profs) == 0 {
		return nil, errors.New("no profile found in profile file")
	}
	return profs, nil
}

func loadProfile(cmd *cobra.Command) *profile.Profile {
	path, _ := cmd.Flags().GetString("config")
	prof, _ := cmd.Flags().GetString("profile")
	var opts []profile.Option
	if prof != "" || path != "" {
		opts = []profile.Option{profile.FromFile(prof, path), profile.MergeWith(profile.FromEnv())}
	}
	p, err := profile.New(opts...)
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
