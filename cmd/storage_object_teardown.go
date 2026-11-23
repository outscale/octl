/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/octl/pkg/debug"
	"github.com/outscale/octl/pkg/messages"
	"github.com/outscale/octl/pkg/spinner"
	"github.com/outscale/octl/pkg/tree"
	"github.com/outscale/osc-sdk-go/v3/pkg/oos"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var objectTeardownCmd = &cobra.Command{
	Use:   "teardown object_key",
	Short: "Tears down an object",
	Args:  cobra.ExactArgs(1),
	Run:   teardownObject,
}

func teardownObject(cmd *cobra.Command, args []string) {
	debug.Println(cmd.Name() + " called")

	if yes, _ := cmd.Flags().GetBool("yes"); !yes {
		displayed := displayObject(cmd, args)
		if !displayed {
			messages.Info("Too many versions to display, skipping tree rendering...")
		}
		if !messages.Prompt("Are you sure you want to teardown this object ?") {
			return
		}
	}

	doObjectTeardown(cmd, args)
}

func displayObject(cmd *cobra.Command, args []string) bool {
	objKey := args[0]
	bucket, _ := cmd.Flags().GetString("bucket")

	ctx := cmd.Context()

	p := loadProfile(cmd)
	cl, err := oos.NewClient(ctx, p, awsOptions(cmd)...)
	if err != nil {
		messages.ExitErr(err)
	}

	paginator := s3.NewListObjectVersionsPaginator(cl, &s3.ListObjectVersionsInput{
		Prefix: &objKey,
		Bucket: &bucket,
	})
	objects, err := paginator.NextPage(ctx)
	if err != nil {
		messages.ExitErr(fmt.Errorf("list bucket object versions: %w", err))
	}

	versions := lo.Filter(objects.Versions, func(v s3types.ObjectVersion, _ int) bool {
		return ptr.From(v.Key) == objKey
	})
	markers := lo.Filter(objects.DeleteMarkers, func(m s3types.DeleteMarkerEntry, _ int) bool {
		return ptr.From(m.Key) == objKey
	})
	if len(versions) == 0 && len(markers) == 0 {
		messages.ExitErr(errors.New("object not found"))
	}

	cancel := spinner.Run(ctx, "Building dependency tree...", false)
	defer cancel()

	root := &Resource{
		Type: Object,
		ID:   objKey,
	}
	var prefix string
	if ptr.From(objects.IsTruncated) {
		prefix = "at least "
	}

	nullVersion := len(versions) == 1 && (versions[0].VersionId == nil || *versions[0].VersionId == "null")
	if !nullVersion {
		suffix := ""
		if len(versions) > 1 {
			suffix = "s"
		}

		root.children = append(root.children, &Resource{
			Type: Text,
			ID:   fmt.Sprintf("%s%d object version%s", prefix, len(versions), suffix),
		})
	}
	if len(markers) > 0 {
		suffix := ""
		if len(markers) > 1 {
			suffix = "s"
		}

		root.children = append(root.children, &Resource{
			Type: Text,
			ID:   fmt.Sprintf("%s%d delete marker%s", prefix, len(markers), suffix),
		})
	}

	err = tree.WriteTo(root, os.Stdout)
	if err != nil {
		messages.ExitErr(err)
	}

	return true
}

func doObjectTeardown(cmd *cobra.Command, args []string) {
	objKey := args[0]
	bucket, _ := cmd.Flags().GetString("bucket")

	ctx := cmd.Context()

	p := loadProfile(cmd)
	cl, err := oos.NewClient(ctx, p, awsOptions(cmd)...)
	if err != nil {
		messages.ExitErr(err)
	}

	ch := make(chan []s3types.ObjectIdentifier)
	group, ctxGroup := errgroup.WithContext(ctx)
	group.Go(func() error {
		return getObjects(ctxGroup, cl, ch, bucket, objKey)
	})
	group.Go(func() error {
		return deleteObjects(ctxGroup, cl, ch, bucket)
	})
	err = group.Wait()
	if err != nil {
		messages.ExitErr(err)
	}

	messages.Success("Object teardown complete")
}
