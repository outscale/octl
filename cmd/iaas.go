/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package cmd

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	commandbuilder "github.com/outscale/octl/pkg/builder/command"
	"github.com/outscale/octl/pkg/config"
	"github.com/outscale/octl/pkg/debug"
	"github.com/outscale/octl/pkg/messages"
	"github.com/outscale/octl/pkg/runner"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/spf13/cobra"
)

// iaasCmd represents the iaascommand
var iaasCmd = &cobra.Command{
	GroupID:           "services",
	Use:               "iaas",
	Short:             "OUTSCALE IaaS management",
	PersistentPreRunE: securityGroupRulePorts,
}

func init() {
	rootCmd.AddCommand(iaasCmd)
	b := commandbuilder.NewBuilder("iaas", "https://docs.outscale.com/api.html")
	b.Build(iaasCmd, oapi)

	cmd, _, err := iaasCmd.Find([]string{"net"})
	if err != nil {
		panic(err)
	}
	cmd.AddCommand(teardownCmd)
	teardownCmd.Flags().Duration("timeout", 10*time.Minute, "Timeout for a single resource deletion")
	teardownCmd.Flags().Bool("teardown-vms", false, "Tears down VM in net")
	cmd.AddCommand(depsCmd)
}

func oapi(cmd *cobra.Command, args []string) {
	p := loadProfile(cmd)
	cl, err := osc.NewClient(p, sdkOptions(cmd)...)
	if err == nil {
		err = deviceToVolumeID(cmd, cl)
	}
	if err == nil {
		err = runner.Run[*osc.Client, *osc.ErrorResponse](cmd, args, cl, config.For("iaas"))
	}
	if err != nil {
		messages.ExitErr(err)
	}
}

var rePorts = regexp.MustCompile("([a-z0-9-]+)(/(-?[0-9]+)(-([0-9]+))?)?")

// securityGroupRulePorts parses the --ports flag into --protocol/--from-port/--to-port flags.
func securityGroupRulePorts(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()
	ports, err := flags.GetStringSlice("ports")
	if err != nil {
		return nil //nolint
	}
	debug.Println("found ports", ports)
	var protocols, fromPorts, toPorts []string
	for _, port := range ports {
		ms := rePorts.FindAllStringSubmatch(port, 1)
		if len(ms) == 0 {
			return nil
		}
		protocol, from, to := ms[0][1], ms[0][3], ms[0][5]
		fromPort, toPort := "-1", "-1"
		if from != "" {
			fromPort = from
		}
		if to != "" {
			toPort = to
		}
		if fromPort != "-1" && toPort == "-1" {
			toPort = fromPort
		}
		if err != nil {
			return fmt.Errorf("invalid ports: %w", err)
		}
		protocols = append(protocols, protocol)
		fromPorts = append(fromPorts, fromPort)
		toPorts = append(toPorts, toPort)
	}
	if f := flags.Lookup("protocol"); f != nil {
		debug.Println("set protocol", protocols)
		_ = setFlag(f, strings.Join(protocols, ","))
	}
	if f := flags.Lookup("from-port"); f != nil {
		debug.Println("set from", fromPorts)
		_ = setFlag(f, strings.Join(fromPorts, ","))
	}
	if f := flags.Lookup("to-port"); f != nil {
		debug.Println("set to", toPorts)
		_ = setFlag(f, strings.Join(toPorts, ","))
	}
	return nil
}
