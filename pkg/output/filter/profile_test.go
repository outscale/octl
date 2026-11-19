/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package filter_test

import (
	"slices"
	"testing"

	"github.com/outscale/octl/pkg/output/filter"
	"github.com/outscale/octl/pkg/output/result"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfileFilter(t *testing.T) {
	seq := func(yield func(result.Result) bool) {
		yield(result.Result{Ok: map[string]any{"VmId": "i-1"}, Profile: "dev"})
		yield(result.Result{Ok: "not a map", Profile: "dev"})
		yield(result.Error(assert.AnError))
	}
	res := slices.Collect(filter.Profile{}.Filter(t.Context(), seq))
	require.Len(t, res, 3)
	assert.Equal(t, "dev", res[0].Ok.(map[string]any)["_profile"])
	assert.Equal(t, "not a map", res[1].Ok)
	assert.Error(t, res[2].Error)
}
