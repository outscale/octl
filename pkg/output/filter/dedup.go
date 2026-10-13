/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package filter

import (
	"context"
	"crypto/sha1" //nolint: gosec
	"encoding/json"
	"iter"
	"slices"

	"github.com/outscale/octl/pkg/output/result"
)

const (
	shaSize     = sha1.Size
	dedupBuffer = 100
)

type Dedup struct {
	resetOnNew bool

	iter    int
	force   bool
	last    [][shaSize]byte
	current [][shaSize]byte
}

func (j *Dedup) Filter(ctx context.Context, seq iter.Seq[result.Result]) iter.Seq[result.Result] {
	return func(yield func(result.Result) bool) {
		for v := range seq {
			if v.Error != nil {
				_ = yield(v)
				return
			}
			if v.Iter != j.iter {
				j.iter = v.Iter
				j.force = false
				j.current, j.last = j.last, j.current
				j.current = j.current[:0]
			}
			js, err := json.Marshal(v.Ok)
			if err == nil {
				hash := sha1.Sum(js) //nolint: gosec
				j.current = append(j.current, hash)
				if !j.force && slices.ContainsFunc(j.last, func(e [shaSize]byte) bool {
					return hash == e
				}) {
					continue
				}
				j.force = j.resetOnNew
			}
			if !yield(v) {
				return
			}
		}
	}
}

func NewDedup(resetOnNew bool) *Dedup {
	return &Dedup{
		resetOnNew: resetOnNew,
		last:       make([][shaSize]byte, dedupBuffer),
		current:    make([][shaSize]byte, dedupBuffer),
	}
}

var _ Interface = (*Dedup)(nil)
