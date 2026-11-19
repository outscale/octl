/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package cmd

import (
	"path/filepath"
	"testing"

	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cmdWithConfig(t *testing.T, file string) *cobra.Command {
	t.Helper()
	cmd := &cobra.Command{}
	cmd.Flags().String("config", "", "")
	cmd.Flags().String("profile", "", "")
	require.NoError(t, cmd.Flags().Set("config", file))
	return cmd
}

func TestLoadProfiles(t *testing.T) {
	file := filepath.Join(t.TempDir(), "config.json")
	cf := &profile.ConfigFile{
		Path: file,
		Profiles: map[string]profile.Profile{
			"b": {AccessKey: "AKB", SecretKey: "SKB"},
			"a": {AccessKey: "AKA", SecretKey: "SKA"},
		},
	}
	require.NoError(t, cf.Save())

	t.Run("all selects every profile sorted by name", func(t *testing.T) {
		cmd := cmdWithConfig(t, file)
		require.NoError(t, cmd.Flags().Set("profile", "all"))
		profs, err := loadProfiles(cmd)
		require.NoError(t, err)
		require.Len(t, profs, 2)
		assert.Equal(t, "a", profs[0].name)
		assert.Equal(t, "AKA", profs[0].p.AccessKey)
		assert.Equal(t, "b", profs[1].name)
	})

	t.Run("a comma separated list selects the listed profiles", func(t *testing.T) {
		cmd := cmdWithConfig(t, file)
		require.NoError(t, cmd.Flags().Set("profile", "b,a"))
		profs, err := loadProfiles(cmd)
		require.NoError(t, err)
		require.Len(t, profs, 2)
		assert.Equal(t, "b", profs[0].name)
		assert.Equal(t, "a", profs[1].name)
	})

	t.Run("an unknown profile returns an error", func(t *testing.T) {
		cmd := cmdWithConfig(t, file)
		require.NoError(t, cmd.Flags().Set("profile", "a,c"))
		_, err := loadProfiles(cmd)
		require.Error(t, err)
	})

	t.Run("a single name keeps the single profile behavior", func(t *testing.T) {
		cmd := cmdWithConfig(t, file)
		require.NoError(t, cmd.Flags().Set("profile", "a"))
		profs, err := loadProfiles(cmd)
		require.NoError(t, err)
		require.Len(t, profs, 1)
		assert.Equal(t, "a", profs[0].name)
		assert.Equal(t, "AKA", profs[0].p.AccessKey)
	})
}
