//go:build debug

/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package debug

import (
	"fmt"
	"os"
)

func Println(s ...any) {
	fmt.Fprintln(os.Stderr, s...)
}
