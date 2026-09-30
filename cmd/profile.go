/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/charmbracelet/huh"
	"github.com/outscale/goutils/sdk/sanitize"
	"github.com/outscale/octl/pkg/config"
	"github.com/outscale/octl/pkg/messages"
	"github.com/outscale/octl/pkg/output"
	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var sanitizer = sanitize.NewStruct(sanitize.MatchField(sanitize.Sensitive), sanitize.RedactField(sanitize.KeepFirst2Last2))

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	GroupID: "config",
	Use:     "profile",
	Short:   "Profile file management",
	Long:    `Creates, updates profile from a config file`,
}

var profileListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Lists all profiles from a config file",
	Aliases: []string{"ls"},
	Run:     listProfiles,
}

var profileCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Display the profile used based on flags/env",
	Run:   currentProfile,
}

var profileAddCmd = &cobra.Command{
	Use:   "add name",
	Short: "Add a profile to a config file",
	Args:  cobra.ExactArgs(1),
	Run:   addProfile,
}

var profileUseCmd = &cobra.Command{
	Use:   "use name",
	Short: "Mark a profile as the default one",
	Args:  cobra.ExactArgs(1),
	Run:   setDefaultProfile,
}

var profileDeleteCmd = &cobra.Command{
	Use:     "delete name",
	Aliases: []string{"del", "rm"},
	Short:   "Delete a profile from a config file",
	Run:     deleteProfile,
}

var profileExplainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Explain current configuration",
	Run:   explainProfile,
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileAddCmd)
	profileCmd.AddCommand(profileUseCmd)
	profileCmd.AddCommand(profileCurrentCmd)
	profileCmd.AddCommand(profileDeleteCmd)
	profileCmd.AddCommand(profileExplainCmd)

	profileAddCmd.Flags().String("ak", "", "Access Key")
	profileAddCmd.Flags().String("sk", "", "Secret Key - if not specified, you prompted to enter it")
	profileAddCmd.Flags().String("region", "eu-west-2", "Region")
	profileAddCmd.Flags().Bool("default", false, "Sets the new profile as the default")
	_ = cobra.MarkFlagRequired(profileAddCmd.Flags(), "ak")
	_ = profileAddCmd.RegisterFlagCompletionFunc(
		"region",
		cobra.FixedCompletions([]cobra.Completion{"eu-west-2", "us-west-1", "us-east-2", "cloudgouv-eu-west-1", "ap-northeast-1"}, cobra.ShellCompDirectiveDefault),
	)
}

type profileEntry struct {
	Name           string `json:"name"`
	Default        bool   `json:"default,omitempty"`
	profile.Fields `json:",inline"`
}

var profileColumns = config.Columns{{Title: "Name", Content: ".name"}, {Title: "Region", Content: ".region"}, {Title: "Default", Content: ".default"}}

func optionsFromCommand(cmd *cobra.Command) profile.Options {
	var opt profile.Options

	path, _ := cmd.Flags().GetString("config")
	if path != "" {
		opt.FilePath = &path
	}

	return opt
}

func loadConfig(cmd *cobra.Command) (*profile.ConfigFile, error) {
	return profile.LoadConfigFile(optionsFromCommand(cmd))
}

func listProfiles(cmd *cobra.Command, _ []string) {
	cf, err := loadConfig(cmd)
	if err != nil {
		messages.ExitErr(err)
	}
	out, _, err := output.NewFromFlags(cmd.Flags(), "table", "", profileColumns, false, true)
	if err != nil {
		messages.ExitErr(err)
	}
	def, _, err := cf.DefaultProfile()
	switch {
	case errors.Is(err, profile.ErrProfileNotFound):
	case err != nil:
		messages.ExitErr(err)
	}
	lst := lo.Map(cf.ProfileList(), func(k string, index int) profileEntry {
		f, _ := cf.Profile(k) // TODO: can fail on malformated profile
		return profileEntry{Name: k, Fields: f, Default: k == def}
	})
	_ = out.Format(cmd.Context(), os.Stdout, sanitizer.Sanitize(lst))
}

func currentProfile(cmd *cobra.Command, _ []string) {
	out, _, err := output.NewFromFlags(cmd.Flags(), "yaml", "", profileColumns, false, true)
	if err != nil {
		messages.ExitErr(err)
	}

	prof := loadProfile(cmd)

	// TODO: export profile status
	_ = out.Format(cmd.Context(), os.Stdout,
		sanitizer.Sanitize(profileEntry{Name: prof.ProfileName, Fields: prof.Values, Default: false}))
}

func addProfile(cmd *cobra.Command, args []string) {
	name := args[0]

	cf, err := loadConfig(cmd)
	if err != nil {
		messages.ExitErr(err)
	}

	var fields profile.Fields
	fields.AccessKey, _ = cmd.Flags().GetString("ak")
	if fields.AccessKey == "" {
		messages.Exit(1, "Access key is required")
	}

	fields.SecretKey, _ = cmd.Flags().GetString("sk")
	if fields.SecretKey == "" {
		var err error
		fields.SecretKey, err = Prompt("Enter the secret Key:", huh.EchoModePassword)
		if err != nil {
			messages.ExitErr(err)
		}
	}
	if fields.SecretKey == "" {
		messages.Exit(1, "Secret key is required")
	}

	fields.Region, _ = cmd.Flags().GetString("region")
	if fields.Region == "" {
		messages.Exit(1, "Region is required")
	}

	err = cf.ProfileAdd(name, fields)
	if err != nil {
		messages.ExitErr(err)
	}

	def, _ := cmd.Flags().GetBool("default")
	if def {
		_ = cf.SetDefault(name)
	}
	err = cf.Save()
	if err != nil {
		messages.ExitErr(err)
	}
	messages.Success("Profile %q has been added", name)
}

func setDefaultProfile(cmd *cobra.Command, args []string) {
	name := args[0]
	cf, err := loadConfig(cmd)
	if err != nil {
		messages.ExitErr(err)
	}
	err = cf.SetDefault(name)
	if err == nil {
		err = cf.Save()
	}
	if err != nil {
		messages.ExitErr(err)
	}
	messages.Success("Profile %q is now the default", name)
}

func deleteProfile(cmd *cobra.Command, args []string) {
	name := args[0]
	cf, err := loadConfig(cmd)
	if err != nil {
		messages.ExitErr(err)
	}

	err = cf.ProfileRemove(name)
	if err != nil {
		err = cf.Save()
	}
	if err != nil {
		messages.ExitErr(err)
	}
	messages.Success("Profile %q has been deleted", name)
}

var explainColumns = config.Columns{{Title: "Entry", Content: ".entry"}, {Title: "Source", Content: ".source"}, {Title: "Value", Content: ".value"}}

type explainEntry struct {
	Entry  string `json:"entry"`
	Source string `json:"source,omitempty"`
	Value  string `json:"value"`
}

func explainProfile(cmd *cobra.Command, _ []string) {
	prof := loadProfile(cmd)
	out, _, err := output.NewFromFlags(cmd.Flags(), "table", "", explainColumns, false, true)
	if err != nil {
		messages.ExitErr(err)
	}

	var entry []explainEntry
	value := reflect.ValueOf(prof.Values)
	kind := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if !field.IsValid() || !field.CanInterface() {
			continue
		}

		fieldType := kind.Field(i)
		source := prof.Sources[fieldType.Name]
		if field.Kind() == reflect.Map {
			iter := field.MapRange()
			for iter.Next() {
				entryName := fieldType.Name + "." + iter.Key().String()
				source := prof.Sources[entryName]
				entry = append(entry, explainEntry{Entry: entryName, Source: source, Value: fmt.Sprint(iter.Value().Interface())})
			}
		} else {
			entry = append(entry, explainEntry{Entry: fieldType.Name, Source: source, Value: fmt.Sprint(field.Interface())})
		}

	}

	_ = out.Format(cmd.Context(), os.Stdout, sanitizer.Sanitize(entry))
}
