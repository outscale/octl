/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package cmd

import (
	"fmt"

	"github.com/outscale/gli/pkg/version"
	sdkversion "github.com/outscale/osc-sdk-go/v3/pkg/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display GLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gli version %s\n", version.Version)
		fmt.Printf("based on Go SDK %s\n", sdkversion.SDKVersion)
		fmt.Printf("based on Outscale API %s\n", sdkversion.APIVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
