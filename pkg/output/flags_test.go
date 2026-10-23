/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package output_test

import (
	"bytes"
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/outscale/octl/pkg/config"
	"github.com/outscale/octl/pkg/output"
	"github.com/outscale/octl/pkg/output/read"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

type mpReq struct{}
type mpItem struct{ Id string }
type mpResp struct{ Items []mpItem }

func mpFetch(profile string, items ...mpItem) read.FetchPage {
	return read.FetchPage{
		Profile: profile,
		Method: reflect.ValueOf(func(context.Context, mpReq) (mpResp, error) {
			return mpResp{Items: items}, nil
		}),
		Args: []reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf(mpReq{})},
	}
}

func TestMultiProfileOutput(t *testing.T) {
	newFlags := func(profile string) *pflag.FlagSet {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		fs.String("profile", "", "")
		fs.String("output", "", "")
		fs.Int("max-pages", 20, "")
		require.NoError(t, fs.Set("profile", profile))
		return fs
	}

	t.Run("table gets a Profile column", func(t *testing.T) {
		fs := newFlags("all")
		_, out, err := output.NewFromFlags(fs, "table", "Items", config.Columns{{Title: "ID", Content: ".Id"}}, false, true)
		require.NoError(t, err)
		buf := &bytes.Buffer{}
		output.InjectOutput(buf)
		defer output.InjectOutput(os.Stdout)
		require.NoError(t, out.Output(t.Context(), mpFetch("a", mpItem{"a1"}), mpFetch("b", mpItem{"b1"})))
		require.Contains(t, buf.String(), "Profile")
		require.Contains(t, buf.String(), "a1")
		require.Contains(t, buf.String(), "b1")
	})

	t.Run("json gets a _profile field", func(t *testing.T) {
		fs := newFlags("a,b")
		_, out, err := output.NewFromFlags(fs, "json", "Items", nil, false, false)
		require.NoError(t, err)
		buf := &bytes.Buffer{}
		output.InjectOutput(buf)
		defer output.InjectOutput(os.Stdout)
		require.NoError(t, out.Output(t.Context(), mpFetch("a", mpItem{"a1"})))
		require.Contains(t, buf.String(), `"_profile": "a"`)
	})

	t.Run("a single profile keeps the standard output", func(t *testing.T) {
		fs := newFlags("a")
		_, out, err := output.NewFromFlags(fs, "table", "Items", config.Columns{{Title: "ID", Content: ".Id"}}, false, true)
		require.NoError(t, err)
		buf := &bytes.Buffer{}
		output.InjectOutput(buf)
		defer output.InjectOutput(os.Stdout)
		require.NoError(t, out.Output(t.Context(), mpFetch("a", mpItem{"a1"})))
		require.NotContains(t, buf.String(), "Profile")
	})
}
