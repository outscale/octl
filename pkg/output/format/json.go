/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>
SPDX-License-Identifier: BSD-3-Clause
*/
package format

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/tidwall/pretty"
)

type JSON struct{}

func (JSON) Format(ctx context.Context, w io.Writer, v any) error {
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}
	if IsTerminal(w) {
		buf = pretty.Color(buf, nil)
	}
	_, err = fmt.Fprintln(w, string(buf))
	return err
}

func (j JSON) Error(ctx context.Context, v any) error {
	return j.Format(ctx, os.Stderr, v)
}

var _ Interface = JSON{}
