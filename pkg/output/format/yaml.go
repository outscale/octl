/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>
SPDX-License-Identifier: BSD-3-Clause
*/
package format

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/goccy/go-yaml"
)

type YAML struct {
	style string
	raw   bool
}

func NewYAML(style string, raw bool) YAML {
	return YAML{style: style, raw: raw}
}

func (y YAML) Format(ctx context.Context, w io.Writer, v any) error {
	buf := new(bytes.Buffer)
	opts := []yaml.EncodeOption{yaml.UseSingleQuote(true), yaml.Indent(2), yaml.CustomMarshaler(
		func(v []byte) ([]byte, error) {
			return []byte(base64.StdEncoding.EncodeToString(v)), nil
		},
	)}
	if !y.raw {
		opts = append(opts, yaml.OmitZero())
	}
	enc := yaml.NewEncoder(buf, opts...)
	err := enc.EncodeContext(ctx, v)
	if err == nil {
		err = enc.Close()
	}
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}
	if IsTerminal(w) && y.style != "" {
		err = quick.Highlight(w, buf.String(), "YAML", "terminal256", y.style)
	} else {
		_, err = fmt.Fprintln(w, buf.String())
	}
	return err
}

func (y YAML) Error(ctx context.Context, v any) error {
	return y.Format(ctx, os.Stderr, v)
}

var _ Interface = YAML{}
