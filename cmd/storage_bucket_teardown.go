/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/charmbracelet/x/term"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/octl/pkg/debug"
	"github.com/outscale/octl/pkg/messages"
	"github.com/outscale/octl/pkg/spinner"
	"github.com/outscale/octl/pkg/tree"
	"github.com/outscale/osc-sdk-go/v3/pkg/oos"
	"github.com/spf13/cobra"
)

var bucketTeardownCmd = &cobra.Command{
	Use:   "teardown bucket_name",
	Short: "Tears down a bucket and its objects",
	Args:  cobra.ExactArgs(1),
	Run:   teardownBucket,
}

const (
	Bucket        ResourceType = "bucket"
	Object        ResourceType = "object"
	ObjectVersion ResourceType = "object version"

	objectDeleteBatchSize int32 = 1000
)

func teardownBucket(cmd *cobra.Command, args []string) {
	debug.Println(cmd.Name() + " called")

	if yes, _ := cmd.Flags().GetBool("yes"); !yes {
		truncated := displayBucket(cmd, args)
		if truncated {
			messages.Info("Too many objects in the bucket, partial tree displayed...")
		}
		if !messages.Prompt("Are you sure you want to teardown this bucket ?") {
			return
		}
	}

	doBucketTeardown(cmd, args)
}

func doBucketTeardown(cmd *cobra.Command, args []string) {
	bucket := args[0]
	ctx := cmd.Context()

	p := loadProfile(cmd)
	cl, err := oos.NewClient(ctx, p, awsOptions(cmd)...)
	if err != nil {
		messages.ExitErr(err)
	}

	ch := make(chan []s3types.ObjectIdentifier)
	group, ctxGroup := errgroup.WithContext(ctx)
	group.Go(func() error {
		return getObjects(ctxGroup, cl, ch, bucket)
	})
	group.Go(func() error {
		return deleteObjects(ctxGroup, cl, ch, bucket)
	})
	err = group.Wait()
	if err != nil {
		messages.ExitErr(err)
	}

	_, err = cl.DeleteBucket(ctx, &s3.DeleteBucketInput{Bucket: &bucket})
	if err != nil {
		messages.ExitErr(err)
	}

	messages.Success("Bucket teardown complete")
}

func getObjects(ctx context.Context, cl *oos.Client, batches chan<- []s3types.ObjectIdentifier, bucket string) error {
	defer close(batches)

	input := &s3.ListObjectVersionsInput{
		Bucket:  &bucket,
		MaxKeys: new(objectDeleteBatchSize),
	}
	paginator := s3.NewListObjectVersionsPaginator(cl, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("list bucket object versions: %w", err)
		}

		objects := make([]s3types.ObjectIdentifier, 0, len(page.Versions)+len(page.DeleteMarkers))
		for _, object := range page.Versions {
			objects = append(objects, s3types.ObjectIdentifier{
				Key:       object.Key,
				VersionId: object.VersionId,
			})
		}
		for _, marker := range page.DeleteMarkers {
			objects = append(objects, s3types.ObjectIdentifier{
				Key:       marker.Key,
				VersionId: marker.VersionId,
			})
		}
		if len(objects) == 0 {
			continue
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case batches <- objects:
		}
	}

	return nil
}

func deleteObjects(ctx context.Context, cl *oos.Client, ch <-chan []s3types.ObjectIdentifier, bucket string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case batch, ok := <-ch:
			if !ok {
				return nil
			}
			if len(batch) == 0 {
				continue
			}

			spinCancel := spinner.Run(ctx, "Deleting batch of "+strconv.Itoa(len(batch))+" objects/versions/markers...", false)
			output, err := cl.DeleteObjects(ctx, &s3.DeleteObjectsInput{
				Bucket: &bucket,
				Delete: &s3types.Delete{
					Objects: batch,
				},
			})
			spinCancel()
			if err != nil {
				return fmt.Errorf("unable to delete bucket objects: %w", err)
			}
			if output == nil {
				return errors.New("unable to delete bucket objects: empty response")
			}

			deleteErrors := lo.Map(output.Errors, func(objectErr s3types.Error, _ int) error {
				return fmt.Errorf(
					"unable to delete object %q version %q: %s: %s",
					ptr.From(objectErr.Key),
					ptr.From(objectErr.VersionId),
					ptr.From(objectErr.Code),
					ptr.From(objectErr.Message),
				)
			})
			err = errors.Join(deleteErrors...)
			if err != nil {
				return err
			}
		}
	}
}

func displayBucket(cmd *cobra.Command, args []string) bool {
	bucketName := args[0]

	ctx := cmd.Context()

	p := loadProfile(cmd)
	cl, err := oos.NewClient(ctx, p, awsOptions(cmd)...)
	if err != nil {
		messages.ExitErr(err)
	}

	_, height, _ := term.GetSize(os.Stdout.Fd())
	displayLimit := height * 2 / 3
	// one line for bucket, one for ellipsis, and 2 lines for each object (object + versions)
	objectLimit := max(10, (displayLimit-2)/2)

	paginator := s3.NewListObjectVersionsPaginator(cl, &s3.ListObjectVersionsInput{
		Bucket: &bucketName,
	})
	objects, err := paginator.NextPage(ctx)
	if err != nil {
		messages.ExitErr(fmt.Errorf("list bucket object versions: %w", err))
	}

	cancel := spinner.Run(ctx, "Building dependency tree...", false)
	defer cancel()

	root := &Resource{
		Type: Bucket,
		ID:   bucketName,
	}

	groups := lo.PartitionBy(objects.Versions, func(obj s3types.ObjectVersion) string {
		return ptr.From(obj.Key)
	})
	pageTruncated := ptr.From(objects.IsTruncated)
	truncated := pageTruncated || len(groups) > objectLimit

	for i, versions := range groups[:min(len(groups), objectLimit)] {
		var prefix, suffix string
		if pageTruncated && i == len(groups)-1 {
			prefix = "at least "
		}
		if len(versions) > 1 {
			suffix = "s"
		}

		obj := &Resource{
			Type: Object,
			ID:   ptr.From(versions[0].Key),
		}

		nullVersion := len(versions) == 1 && (versions[0].VersionId == nil || *versions[0].VersionId == "null")
		if !nullVersion {
			obj.children = []*Resource{
				{
					Type: Text,
					ID:   fmt.Sprintf("%s%d object version%s", prefix, len(versions), suffix),
				},
			}
		}

		root.children = append(root.children, obj)
	}

	if truncated {
		root.children = append(root.children, &Resource{
			Type: Text,
			ID:   "...",
		})
	}

	err = tree.WriteTo(root, os.Stdout)
	if err != nil {
		messages.ExitErr(err)
	}

	return truncated
}
