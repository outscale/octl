/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package cmd_test

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestBucketTeardown(t *testing.T) {
	t.Run("teardown on an empty bucket", func(t *testing.T) {
		sum := sha1.Sum([]byte(t.TempDir()))
		bucket := hex.EncodeToString(sum[:])
		_ = run(t, []string{"storage", "bucket", "create", "--bucket", bucket}, nil)

		_ = run(t, []string{"storage", "bucket", "teardown", bucket, "-y"}, nil)
		runWithError(t, []string{"storage", "bucket", "describe", bucket}, nil)
	})
	t.Run("teardown on a bucket with versioned objects", func(t *testing.T) {
		sum := sha1.Sum([]byte(t.TempDir()))
		bucket := hex.EncodeToString(sum[:])

		_ = run(t, []string{"storage", "bucket", "create", "--bucket", bucket}, nil)
		_ = run(t, []string{"storage", "bucket", "versioning", "enable", bucket}, nil)

		for _, object := range []string{"first.txt", "second.txt"} {
			for v := 1; v <= 2; v++ {
				path := file(t, fmt.Sprintf("%s-%d", object, v), fmt.Sprintf("version %d", v))

				_ = run(t, []string{"storage", "object", "put", object, "--bucket", bucket, "--body", path}, nil)
			}
		}

		_ = run(t, []string{"storage", "bucket", "teardown", bucket, "-y"}, nil)

		runWithError(t, []string{"storage", "bucket", "describe", bucket}, nil)
	})
}
