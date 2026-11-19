/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package filter

import (
	"context"
	"iter"

	"github.com/outscale/octl/pkg/output/result"
)

// Profile injects the source profile of a result as a _profile field.
type Profile struct{}

func (Profile) Filter(_ context.Context, seq iter.Seq[result.Result]) iter.Seq[result.Result] {
	return func(yield func(result.Result) bool) {
		for v := range seq {
			if v.Error == nil {
				if content, ok := v.Ok.(map[string]any); ok {
					content["_profile"] = v.Profile
				}
			}
			if !yield(v) {
				return
			}
		}
	}
}

var _ Interface = Profile{}
