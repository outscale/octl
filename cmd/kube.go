/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package cmd

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/outscale/octl/pkg/alias"
	"github.com/outscale/octl/pkg/builder"
	"github.com/outscale/octl/pkg/config"
	"github.com/outscale/octl/pkg/debug"
	"github.com/outscale/octl/pkg/messages"
	"github.com/outscale/octl/pkg/preferences"
	"github.com/outscale/octl/pkg/runner"
	"github.com/outscale/osc-sdk-go/v3/pkg/oks"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// oksCmd represents the kubecommand
var oksCmd = &cobra.Command{
	GroupID: "services",
	Use:     "kube",
	Short:   "OUTSCALE Kubernetes as a Service (OKS) management",
	Aliases: []string{"oks"},
}

var projectUseCmd = &cobra.Command{
	Use:   "use [project_id_or_name]",
	Short: "Set a default project for cluster commands, reset it without args",
	Run:   useProject,
	Args:  cobra.MaximumNArgs(1),
}

func init() {
	rootCmd.AddCommand(oksCmd)
	b := builder.NewBuilder[oks.Client]("kube", "https://docs.outscale.com/oks.html")
	b.BuildAPI(oksCmd, func(m reflect.Method) bool {
		return m.Type.NumIn() >= 3 && m.Type.NumOut() == 2 && !strings.HasSuffix(m.Name, "Raw") &&
			!strings.HasSuffix(m.Name, "WithBody")
	}, kube)
	b.Build(oksCmd, nil)

	// Add --project flag to kube cluster commands
	clusterCmd, _ := lo.Find(oksCmd.Commands(), func(c *cobra.Command) bool { return c.Name() == "cluster" })
	for _, cmd := range clusterCmd.Commands() {
		f := cmd.Flags().Lookup("project")
		if f == nil {
			cmd.Flags().String("project", preferences.Preferences.Kube.DefaultProject, "project name")
			alias.SetDefault(cmd.Flags().Lookup("project"), preferences.Preferences.Kube.DefaultProject)
		} else {
			f.DefValue = preferences.Preferences.Kube.DefaultProject
			alias.SetDefault(f, preferences.Preferences.Kube.DefaultProject)
			if _, ok := f.Annotations[cobra.BashCompOneRequiredFlag]; ok && preferences.Preferences.Kube.DefaultProject != "" {
				delete(f.Annotations, cobra.BashCompOneRequiredFlag)
			}
		}
	}

	// Remap <cluster ls --project name> to <project clusters name>
	projectCmd, _ := lo.Find(oksCmd.Commands(), func(c *cobra.Command) bool { return c.Name() == "project" })
	projectClustersCmd, _ := lo.Find(projectCmd.Commands(), func(c *cobra.Command) bool { return c.Name() == "clusters" })
	clusterListCmd, _ := lo.Find(clusterCmd.Commands(), func(c *cobra.Command) bool { return c.Name() == "list" })
	runClusterListCmd := clusterListCmd.Run
	clusterListCmd.Run = func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		if project != "" {
			alias.Reset(cmd.Flag("project"))
			projectClustersCmd.Run(cmd, []string{project})
		} else {
			runClusterListCmd(cmd, nil)
		}
	}

	// Project use
	projectCmd.AddCommand(projectUseCmd)

	// Add --ProjectId flag to kube api cluster commands
	apiCmd, _ := lo.Find(oksCmd.Commands(), func(c *cobra.Command) bool { return c.Name() == "api" })
	for _, cmd := range apiCmd.Commands() {
		switch cmd.Name() {
		case "GetCluster", "UpdateCluster", "DeleteCluster", "GetKubeconfig":
			cmd.Flags().String("ProjectId", preferences.Preferences.Kube.DefaultProject, "project name or id")
		}
	}
}

func kube(cmd *cobra.Command, args []string) {
	debug.Println(cmd.Name() + " called")
	p := loadProfile(cmd)
	cl, err := oks.NewClient(p, sdkOptions(cmd)...)
	if err == nil {
		args = argNameToID(cmd, args, cl)
		flagNameToID(cmd, cl)
		err = runner.Run[*oks.Client, *oks.ErrorResponse](cmd, args, cl, config.For("kube"))
	}
	if err != nil {
		messages.ExitErr(err)
	}
}

func argNameToID(cmd *cobra.Command, args []string, cl *oks.Client) []string {
	if len(args) == 0 {
		debug.Println("no arg to replace")
		return args
	}
	if _, err := uuid.Parse(args[0]); err == nil {
		debug.Println("arg is an uuid")
		return args
	}
	var err error
	switch cmd.Name() {
	case "GetProject", "DeleteProject", "GetProjectNets", "GetProjectQuotas", "GetProjectPublicIps", "GetProjectSnapshots":
		args[0], err = projectNameToID(cmd.Context(), args[0], cl)
	case "GetCluster", "UpdateCluster", "DeleteCluster", "GetKubeconfig":
		project, _ := cmd.Flags().GetString("ProjectId")
		args[0], err = clusterNameToID(cmd.Context(), args[0], project, cl)
	}
	if err != nil {
		messages.ExitErr(err)
	}
	return args
}

func flagNameToID(cmd *cobra.Command, cl *oks.Client) {
	f := cmd.Flags().Lookup("ProjectId")
	if f == nil {
		debug.Println("no flag to replace")
		return
	}
	id, err := projectNameToID(cmd.Context(), f.Value.String(), cl)
	if err != nil {
		messages.ExitErr(err)
	}
	_ = f.Value.Set(id)
}

func projectNameToID(ctx context.Context, name string, cl *oks.Client) (string, error) {
	pjs, err := cl.ListProjects(ctx, &oks.ListProjectsParams{Name: &name})
	if err != nil {
		return "", err
	}
	switch len(pjs.Projects) {
	case 0:
		return "", fmt.Errorf("project %q not found", name)
	default:
		debug.Println("replacing", name, "by", pjs.Projects[0].Id)
		return pjs.Projects[0].Id, nil
	}
}

func clusterNameToID(ctx context.Context, name, project string, cl *oks.Client) (string, error) {
	var (
		cs  *oks.ClusterResponseList
		err error
	)
	if project != "" {
		var pid string
		if _, err := uuid.Parse(project); err == nil {
			pid = project
			debug.Println("project is an uuid")
		} else {
			pid, err = projectNameToID(ctx, project, cl)
			if err != nil {
				return "", err
			}
		}
		cs, err = cl.ListClustersByProjectID(ctx, &oks.ListClustersByProjectIDParams{Name: &name, ProjectId: &pid})
	} else {
		cs, err = cl.ListAllClusters(ctx, &oks.ListAllClustersParams{Name: &name})
	}
	if err != nil {
		return "", err
	}
	clusters := cs.Clusters
	switch len(clusters) {
	case 0:
		return "", fmt.Errorf("cluster %q not found", name)
	case 1:
		debug.Println("replacing", name, "by", clusters[0].Id)
		return clusters[0].Id, nil
	default:
		return "", fmt.Errorf("multiple clusters found with the name %q", name)
	}
}

func useProject(c *cobra.Command, args []string) {
	var def string
	if len(args) > 0 {
		def = args[0]
	}
	preferences.Preferences.Kube.DefaultProject = def
	err := preferences.Preferences.Save()
	if err != nil {
		messages.ExitErr(err)
	}
	if def == "" {
		messages.Success("The default project has been reset")
	} else {
		messages.Success("%q is now the default project for all octl kube commands", def)
	}
}
