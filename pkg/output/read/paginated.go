/*
SPDX-FileCopyrightText: 2025 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package read

import (
	"context"
	"errors"
	"iter"
	"reflect"
	"strings"

	"github.com/outscale/octl/pkg/debug"
	"github.com/outscale/octl/pkg/messages"
	"github.com/outscale/octl/pkg/output/result"
)

type Paginated struct {
	maxPages     int
	contentField string
}

func NewPaginated(contentField string, maxPages int) *Paginated {
	return &Paginated{
		maxPages:     maxPages,
		contentField: contentField,
	}
}

// Read reads the results of one or more calls. Calls are made sequentially,
// in the order of fetches. When a call fails and several fetches are given,
// the error is turned into a warning and the next calls are still read.
func (p *Paginated) Read(ctx context.Context, iter int, fetches ...FetchPage) iter.Seq[result.Result] {
	return func(yield func(result.Result) bool) {
	read:
		for _, fetch := range fetches {
			for res := range p.readOne(ctx, fetch, iter) {
				if res.Error != nil && len(fetches) > 1 {
					messages.Warn("profile %q: %s", fetch.Profile, res.Error)
					continue
				}
				if !yield(res) {
					break read
				}
			}
		}
	}
}

// readOne reads all pages of a single call.
func (p *Paginated) readOne(ctx context.Context, fetch FetchPage, iter int) iter.Seq[result.Result] {
	pager := PagerFor(fetch)
	return func(yield func(result.Result) bool) {
		fetch := fetch
		idx := 0
		for range p.maxPages {
			vres := fetch.Call(ctx)
			if len(vres) == 0 {
				_ = yield(result.Error(errors.New("no result from call")))
				return
			}
			if err, ok := reflect.TypeAssert[error](vres[len(vres)-1]); ok && err != nil {
				_ = yield(result.Error(err))
				return
			}
			if len(vres) < 2 {
				debug.Println("not enough result")
				_ = yield(result.Result{SingleEntry: true, Iter: iter, Profile: fetch.Profile})
				return
			}
			res := reflect.Indirect(vres[0])
			content := res
			if p.contentField != "" && p.contentField != "." {
				content = contentField(content, p.contentField)
			}
			addPreview(content)
			if !content.IsValid() || !content.CanInterface() {
				debug.Println("no content ?")
				_ = yield(result.Result{SingleEntry: true, Iter: iter, Profile: fetch.Profile})
				return
			}
			if content.Kind() != reflect.Slice {
				debug.Println("not a slice", content.Kind())
				_ = yield(result.Result{Ok: content.Interface(), SingleEntry: true, Iter: iter, Profile: fetch.Profile})
				return
			}
			for i := range content.Len() {
				if !yield(result.Result{Ok: content.Index(i).Interface(), Iter: iter, Profile: fetch.Profile}) {
					debug.Println("end")
					return
				}
			}
			idx += content.Len()
			if !pager.HasMore(res) {
				debug.Println("end of list...")
				return
			}
			debug.Println("has more items...")
			var ok bool
			fetch, ok = pager.NextItem(res, fetch, idx)
			if !ok {
				return
			}
		}
	}
}

func contentField(v reflect.Value, path string) reflect.Value {
	before, after, found := strings.Cut(path, ".")
	content := reflect.Indirect(v.FieldByName(before))
	debug.Println("contentField", before, after, content.IsValid())
	if !found {
		return content
	}
	return contentField(content, after)
}

var _ Interface = (*Paginated)(nil)
