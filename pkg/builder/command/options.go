/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package commandbuilder

import (
	"os"
	"strconv"
	"strings"

	"github.com/outscale/octl/pkg/debug"
	"github.com/samber/lo"
)

type numEntriesInSlices map[string]int

var osNumEntriesInSlices = getNumEntriesInSlices(os.Args)

// We parse the command arguments to find index-based flags count the number of flags for each prefix.
// The cobra commands will be build with all the necessary flags (+1 to allow autompletion of next)
func getNumEntriesInSlices(args []string) numEntriesInSlices {
	// count the number of flags
	cnt := lo.CountBy(args, func(arg string) bool {
		return strings.HasPrefix(arg, "--")
	})
	// worst case = 1 index per flag
	num := make(numEntriesInSlices)
	for i := range cnt {
		idxStr := "." + strconv.Itoa(i) + "."
		for _, arg := range args {
			parts := strings.Split(strings.TrimPrefix(arg, "--"), idxStr)
			if len(parts) == 1 {
				continue
			}
			prefix := ""
			for iarg := range len(parts) - 1 {
				num[prefix+parts[iarg]] = i + 1
				prefix += parts[iarg] + idxStr
			}
		}
	}
	debug.Println("numEntriesInSlices", num)
	return num
}

func (n numEntriesInSlices) forPrefix(prefix string) int {
	if n == nil {
		return 1
	}
	if n, found := n[prefix]; found {
		return n + 1
	}
	return 1
}
