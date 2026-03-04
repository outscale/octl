/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>
SPDX-License-Identifier: BSD-3-Clause
*/
package format

import (
	"context"
	"io"
)

type None struct{}

func (None) Format(ctx context.Context, _ io.Writer, _ any) error {
	return nil
}

func (None) Error(ctx context.Context, v any) error {
	return YAML{}.Error(ctx, v)
}

var _ Interface = None{}
