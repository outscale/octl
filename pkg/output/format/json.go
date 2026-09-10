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

	"github.com/alecthomas/chroma/v2/quick"
)

type JSON struct {
	style string
}

func NewJSON(style string) JSON {
	return JSON{style: style}
}

func (j JSON) Format(ctx context.Context, w io.Writer, v any) error {
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}
	if IsTerminal(w) && j.style != "" {
		err = quick.Highlight(w, string(buf), "JSON", "terminal256", j.style)
	} else {
		_, err = fmt.Fprintln(w, string(buf))
	}
	return err
}

func (j JSON) Error(ctx context.Context, v any) error {
	return j.Format(ctx, os.Stderr, v)
}

var _ Interface = JSON{}
