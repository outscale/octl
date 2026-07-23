package cmd

import (
	"reflect"
	"slices"

	oksv1beta2 "github.com/outscale/goutils/oks/clientset/typed/oks.dev/v1beta2"
	"github.com/outscale/octl/pkg/builder"
	"github.com/outscale/octl/pkg/flags"
	"github.com/outscale/octl/pkg/preferences"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// oksCmd represents the kubecommand
var nodepoolCmd = &cobra.Command{
	GroupID: "service",
	Use:     "nodepool",
	Short:   "nodepool commands",
	Aliases: []string{"np"},
}

func init() {
	oksCmd.AddCommand(nodepoolCmd)
	b := builder.NewBuilder[oksv1beta2.NodePoolInterface]("kubeclient_nodepool", "https://docs.outscale.com/api.html")
	b.BuildAPI(nodepoolCmd, func(m reflect.Method) bool {
		return slices.Contains([]string{"List", "Get", "Create", "Update", "Delete"}, m.Name)
	}, kubeapi("kubeclient_nodepool"))
	apiCmd, _ := lo.Find(nodepoolCmd.Commands(), func(c *cobra.Command) bool { return c.Name() == "api" })
	b.Build(oksCmd, apiCmd)
	for _, cmd := range nodepoolCmd.Commands() {
		if cmd.Name() == "api" {
			cmd.PersistentFlags().String("cluster", "", "[REQUIRED] ID of cluster")
			_ = cmd.MarkPersistentFlagRequired("cluster")
		} else {
			cmd.Flags().String("cluster", "", "[REQUIRED] Name or ID of cluster")
			cmd.Flags().String("project", preferences.Preferences.Kube.DefaultProject, "Name or ID of project")
			_ = cmd.MarkFlagRequired("cluster")
			_ = flags.MarkAsNoForward(cmd.Flags(), "project")
		}
	}
}
