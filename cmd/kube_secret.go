/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package cmd

import (
	"os"

	"github.com/outscale/octl/pkg/markdown"
	"github.com/outscale/octl/pkg/output/format"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

var secretLongUse = "Create secret for CCM/CSI driver/Cluster-API provider deployments using the selected AK/SK.\n\n" +
	"CCM:\n" +
	"```shell\n" +
	"octl kube secret --name osc-secret | kubectl apply -f -\n" +
	"```\n" +
	"CSI driver:\n" +
	"```shell\n" +
	"octl kube secret --name osc-csi-bsu | kubectl apply -f -\n" +
	"```\n" +
	"Cluster-API:\n" +
	"```shell\n" +
	"octl kube secret --name cluster-api-provider-outscale --namespace cluster-api-provider-outscale-system | kubectl apply -f -\n" +
	"```\n"

var (
	secretLongUseRendered = lo.Must(markdown.NewRenderer().Render(secretLongUse))
	secretCmd             = &cobra.Command{
		Use:   "secret",
		Long:  secretLongUseRendered,
		Short: "Create secret for CCM or CSI driver deployment",
		Args:  cobra.NoArgs,
		Run:   buildSecret,
	}
)

func init() {
	oksCmd.AddCommand(secretCmd)
	secretCmd.Flags().String("name", "", "[REQUIRED] name of secret")
	secretCmd.Flags().String("namespace", "kube-system", "namespace of secret")
	_ = secretCmd.MarkFlagRequired("name")
}

func buildSecret(cmd *cobra.Command, args []string) {
	p := loadProfile(cmd)
	ns, _ := cmd.Flags().GetString("namespace")
	name, _ := cmd.Flags().GetString("name")
	s := corev1.Secret{
		APIVersion: "v1",
		Kind:       "Secret",
		Name:       name,
		Namespace:  ns,
		Data: map[string][]byte{
			"access_key": []byte(p.AccessKey),
			"secret_key": []byte(p.SecretKey),
			"region":     []byte(p.Region),
		},
		Type: corev1.SecretTypeOpaque,
	}
	_ = format.YAML{}.Format(cmd.Context(), os.Stdout, s)
}
