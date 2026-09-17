/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>
SPDX-License-Identifier: BSD-3-Clause
*/
package format

import (
	"context"
	"crypto/sha1" //nolint: gosec
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type PrintPath struct {
	fmter Interface
}

func NewPrintPath(fmter Interface) PrintPath {
	return PrintPath{fmter: fmter}
}

func (p PrintPath) Format(ctx context.Context, w io.Writer, v any) error {
	dir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("get cache dir: %w", err)
	}
	path := filepath.Join(dir, "octl", "print-path")
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		err = os.MkdirAll(path, 0o700)
		if err != nil {
			return fmt.Errorf("create temp dir: %w", err)
		}
	}
	hash := sha1.Sum([]byte(strings.Join(os.Args, " "))) //nolint: gosec
	s := hex.EncodeToString(hash[:])
	path = filepath.Join(path, s)
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600) //nolint: gosec
	if err == nil {
		err = p.fmter.Format(ctx, fd, v)
	}
	if fd != nil {
		cerr := fd.Close()
		if cerr != nil && err == nil {
			err = cerr
		}
	}
	if err != nil {
		return fmt.Errorf("temp file: %w", err)
	}
	_, err = fmt.Fprint(w, path)
	return err
}

func (p PrintPath) Error(ctx context.Context, v any) error {
	return p.fmter.Error(ctx, v)
}

var _ Interface = PrintPath{}
