/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package messages

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/mattn/go-isatty"
	"github.com/outscale/octl/pkg/style"
)

func Prompt(title string) bool {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		Exit(1, "unable to confirm action, aborting. Add -y to skip confirmation prompt")
	}
	var yes bool
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Affirmative("✅ Yes").
				Negative("❌ No").
				Value(&yes),
		),
	).WithTheme(style.Theme()).Run()
	if err != nil {
		ExitErr(err)
	}
	return yes
}
