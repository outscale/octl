/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package cmd

import (
	"fmt"

	"github.com/outscale/osc-sdk-go/v3/pkg/utils"
	"github.com/spf13/cobra"
)

const Version = "v0.0.0-dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display GLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gli version %s\n", Version)
		fmt.Printf("based on Go SDK %s\n", utils.SdkVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
