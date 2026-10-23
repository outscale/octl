/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package read_test

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/outscale/octl/pkg/output/read"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type multiReq struct{}
type multiResp struct {
	Items []string
}

func fetchFor(profile string, items []string, err error) read.FetchPage {
	return read.FetchPage{
		Profile: profile,
		Method: reflect.ValueOf(func(context.Context, multiReq) (multiResp, error) {
			return multiResp{Items: items}, err
		}),
		Args: []reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf(multiReq{})},
	}
}

func TestMultiProfileRead(t *testing.T) {
	p := read.NewPaginated("Items", 20)

	t.Run("results of all fetches are read in order and annotated with their profile", func(t *testing.T) {
		res := slices.Collect(p.Read(t.Context(), 1,
			fetchFor("a", []string{"a1", "a2"}, nil),
			fetchFor("b", []string{"b1"}, nil),
		))
		require.Len(t, res, 3)
		assert.Equal(t, "a1", res[0].Ok)
		assert.Equal(t, "a", res[0].Profile)
		assert.Equal(t, "a2", res[1].Ok)
		assert.Equal(t, "b1", res[2].Ok)
		assert.Equal(t, "b", res[2].Profile)
	})

	t.Run("a failing fetch is skipped when several fetches are given", func(t *testing.T) {
		res := slices.Collect(p.Read(t.Context(), 1,
			fetchFor("a", nil, errors.New("boom")),
			fetchFor("b", []string{"b1"}, nil),
		))
		require.Len(t, res, 1)
		assert.Equal(t, "b1", res[0].Ok)
		assert.Equal(t, "b", res[0].Profile)
	})

	t.Run("a failing single fetch returns an error", func(t *testing.T) {
		res := slices.Collect(p.Read(t.Context(), 1, fetchFor("a", nil, errors.New("boom"))))
		require.Len(t, res, 1)
		assert.Error(t, res[0].Error)
	})
}

func TestMultiProfileRawRead(t *testing.T) {
	p := read.NewRaw()

	t.Run("a failing fetch is skipped when several fetches are given", func(t *testing.T) {
		res := slices.Collect(p.Read(t.Context(), 1,
			fetchFor("a", nil, errors.New("boom")),
			fetchFor("b", nil, nil),
		))
		require.Len(t, res, 1)
		assert.Equal(t, "b", res[0].Profile)
		assert.Equal(t, multiResp{}, res[0].Ok)
	})

	t.Run("a failing single fetch returns an error", func(t *testing.T) {
		res := slices.Collect(p.Read(t.Context(), 1, fetchFor("a", nil, errors.New("boom"))))
		require.Len(t, res, 1)
		assert.Error(t, res[0].Error)
	})
}
