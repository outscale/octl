/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package cmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/outscale/gli/cmd"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAPI(t *testing.T) {
	t.Run("ReadVms works", func(t *testing.T) {
		os.Args = []string{"gli", "oapi", "ReadVms"}
		stderr, stdout := os.Stderr, os.Stdout
		defer func() {
			os.Stderr, os.Stdout = stderr, stdout
		}()
		dir := t.TempDir()
		var err error
		os.Stderr, err = os.Create(filepath.Join(dir, "stderr")) //nolint
		require.NoError(t, err)
		os.Stdout, err = os.Create(filepath.Join(dir, "stdout")) //nolint
		require.NoError(t, err)
		cmd.Execute()

		err = os.Stderr.Close()
		require.NoError(t, err)
		content, err := os.ReadFile(os.Stderr.Name())
		require.NoError(t, err)
		assert.Empty(t, content)

		err = os.Stdout.Close()
		require.NoError(t, err)
		content, err = os.ReadFile(os.Stdout.Name())
		require.NoError(t, err)
		resp := osc.ReadVmsResponse{}
		err = json.Unmarshal(content, &resp)
		require.NoError(t, err)
		assert.NotEmpty(t, *resp.Vms)
		for _, vm := range *resp.Vms {
			assert.Equal(t, osc.VmStateRunning, vm.State)
		}
		require.NotNil(t, resp.ResponseContext)
		assert.NotEmpty(t, resp.ResponseContext.RequestId)
	})
}
