/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package read

import (
	"context"
	"iter"
	"reflect"

	"github.com/outscale/octl/pkg/messages"
	"github.com/outscale/octl/pkg/output/result"
)

type Raw struct{}

func NewRaw() *Raw {
	return &Raw{}
}

// Read reads the raw response of one or more calls. When a call fails and
// several fetches are given, the error is turned into a warning and the next
// calls are still read.
func (p *Raw) Read(ctx context.Context, iter int, fetches ...FetchPage) iter.Seq[result.Result] {
	return func(yield func(result.Result) bool) {
		for _, fetch := range fetches {
			vres := fetch.Call(ctx)
			if len(vres) == 0 {
				continue
			}
			if err, ok := reflect.TypeAssert[error](vres[len(vres)-1]); ok && err != nil {
				if len(fetches) > 1 {
					messages.Warn("profile %q: %s", fetch.Profile, err)
					continue
				}
				_ = yield(result.Error(err))
				return
			}
			if len(vres) < 2 {
				_ = yield(result.Result{SingleEntry: true, Iter: iter, Profile: fetch.Profile})
				continue
			}
			res := vres[0]
			addPreview(res)
			if !yield(result.Result{Ok: res.Interface(), SingleEntry: true, Iter: iter, Profile: fetch.Profile}) {
				return
			}
		}
	}
}

var _ Interface = (*Raw)(nil)
