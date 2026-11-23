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

func TestObjectTeardown(t *testing.T) {
	t.Run("teardown on a single object", func(t *testing.T) {
		sum := sha1.Sum([]byte(t.TempDir()))
		bucket := hex.EncodeToString(sum[:])

		_ = run(t, []string{"storage", "bucket", "create", "--bucket", bucket}, nil)
		defer func() {
			_ = run(t, []string{"storage", "bucket", "del", bucket, "-y"}, nil)
		}()

		object := "object.txt"
		path := file(t, object, hello)

		_ = run(t, []string{"storage", "object", "put", object, "--bucket", bucket, "--body", path}, nil)

		_ = run(t, []string{"storage", "object", "teardown", object, "--bucket", bucket, "-y"}, nil)

		runWithError(t, []string{"storage", "object", "describe", object, "--bucket", bucket}, nil)
	})
	t.Run("teardown of an object with multiple versions", func(t *testing.T) {
		sum := sha1.Sum([]byte(t.TempDir()))
		bucket := hex.EncodeToString(sum[:])

		_ = run(t, []string{"storage", "bucket", "create", "--bucket", bucket}, nil)
		defer func() {
			_ = run(t, []string{"storage", "bucket", "teardown", bucket, "-y"}, nil)
		}()

		_ = run(t, []string{"storage", "bucket", "versioning", "enable", bucket}, nil)

		keep := "keep.txt"
		tear := "tear.txt"

		_ = run(t, []string{"storage", "object", "put", keep, "--bucket", bucket, "--body", file(t, "obj1", keep)}, nil)

		for v := 1; v <= 2; v++ {
			path := file(t, fmt.Sprintf("%s-%d", tear, v), fmt.Sprintf("version %d", v))

			_ = run(t, []string{"storage", "object", "put", tear, "--bucket", bucket, "--body", path}, nil)
		}

		_ = run(t, []string{"storage", "object", "teardown", tear, "--bucket", bucket, "-y"}, nil)

		_ = run(t, []string{"storage", "object", "describe", keep, "--bucket", bucket}, nil)
		runWithError(t, []string{"storage", "object", "describe", tear, "--bucket", bucket}, nil)
	})
}
