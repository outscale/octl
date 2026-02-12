/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>
SPDX-License-Identifier: BSD-3-Clause
*/
package output

import "context"

type Output interface {
	Content(ctx context.Context, v any) error
	Error(ctx context.Context, v any) error
}
