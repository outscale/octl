/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package options

import (
	"os"
	"strconv"
)

var NumEntriesInSlices = 1

func init() {
	num := os.Getenv("NUM_ENTRIES")
	if num != "" {
		NumEntriesInSlices, _ = strconv.Atoi(num)
	}
}
