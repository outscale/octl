/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package sdk

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/outscale/goutils/sdk/sanitize"
)

type VerboseLogger struct{}

func (VerboseLogger) RequestHttp(ctx context.Context, req *http.Request) {
	req = sanitize.HTTPRequest(req)
	fmt.Fprintf(os.Stderr, "- REQUEST -------------------\n\n%s %s\n\n", req.Method, req.URL)
	_ = req.Header.Write(os.Stderr)
	fmt.Fprintln(os.Stderr)
	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		fmt.Fprintf(os.Stderr, "%s\n\n", string(body))
	}
	fmt.Fprint(os.Stderr, "- REQUEST -------------------\n\n")
}

func (VerboseLogger) ResponseHttp(ctx context.Context, resp *http.Response, d time.Duration) {
	resp = sanitize.HTTPResponse(resp)
	fmt.Fprintf(os.Stderr, "- RESPONSE ------------------\n\n%s\n\n", resp.Status)
	_ = resp.Header.Write(os.Stderr)
	fmt.Fprintln(os.Stderr)
	body, _ := io.ReadAll(resp.Body)
	fmt.Fprintf(os.Stderr, "%s\n\n", string(body))
	_ = resp.Body.Close()
	fmt.Fprint(os.Stderr, "- RESPONSE ------------------\n\n")
}

func (VerboseLogger) Request(ctx context.Context, req any) {}

func (VerboseLogger) Response(ctx context.Context, resp any) {}

func (VerboseLogger) Error(ctx context.Context, err error) {}
