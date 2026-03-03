/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package update

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"

	"github.com/google/go-github/v82/github"
	"github.com/minio/selfupdate"
	"github.com/outscale/octl/pkg/debug"
	"github.com/outscale/octl/pkg/messages"
	"github.com/outscale/octl/pkg/version"
	"github.com/sigstore/sigstore-go/pkg/bundle"
	"github.com/sigstore/sigstore-go/pkg/root"
	"github.com/sigstore/sigstore-go/pkg/tuf"
	"github.com/sigstore/sigstore-go/pkg/verify"
	"golang.org/x/mod/semver"
)

func latestRelease(ctx context.Context) (*github.RepositoryRelease, error) {
	client := github.NewClient(nil)
	rel, _, err := client.Repositories.GetLatestRelease(ctx, "outscale", "octl")
	if err != nil {
		debug.Println("github error", err)
		return nil, nil
	}
	if rel == nil || rel.TagName == nil {
		return nil, nil
	}
	return rel, nil
}

func LatestRelease(ctx context.Context) (string, error) {
	rel, err := latestRelease(ctx)
	if rel == nil || err != nil {
		return "", err
	}
	return rel.GetTagName(), nil
}

func Update(ctx context.Context) error {
	rel, err := latestRelease(ctx)
	if err != nil {
		return err
	}
	if rel == nil {
		return errors.New("no new version found")
	}
	if semver.Compare(version.Version, rel.GetTagName()) >= 0 {
		messages.Warn("Already using the latest version")
		return nil
	}

	debug.Println("OS, arch:", runtime.GOOS, runtime.GOARCH)
	suffix := runtime.GOOS + "_"
	switch runtime.GOARCH {
	case "amd64":
		suffix += "x86_64"
	default:
		suffix += runtime.GOARCH
	}
	if runtime.GOOS == "windows" {
		suffix += ".exe"
	}

	var assetToDownload, checksums, bundle *github.ReleaseAsset
	for _, a := range rel.Assets {
		if a.Name == nil || a.BrowserDownloadURL == nil {
			continue
		}
		if strings.HasSuffix(strings.ToLower(*a.Name), suffix) {
			debug.Println("found asset", a.GetName(), a.GetBrowserDownloadURL())
			assetToDownload = a
		} else if strings.HasSuffix(strings.ToLower(*a.Name), ".sigstore.json") {
			debug.Println("found bundle", a.GetName(), a.GetBrowserDownloadURL())
			bundle = a
		} else if strings.HasSuffix(strings.ToLower(*a.Name), "_checksums.txt") {
			debug.Println("found checksums", a.GetName(), a.GetBrowserDownloadURL())
			checksums = a
		}
	}

	if bundle == nil || checksums == nil || assetToDownload == nil {
		return fmt.Errorf("could not find required assets in release %s", rel.GetTagName())
	}

	b, err := downloadBundle(ctx, bundle)
	if err != nil {
		return err
	}

	cs, err := downloadAndVerifyChecksum(ctx, checksums, b)
	if err != nil {
		return err
	}

	digest, err := findAndCheckAssetDigest(assetToDownload, cs)
	if err != nil {
		return err
	}

	if err := update(ctx, rel.GetTagName(), assetToDownload, digest); err != nil {
		return err
	}

	return nil
}

func findAndCheckAssetDigest(a *github.ReleaseAsset, cs *map[string]string) (string, error) {
	var digest string
	for k, v := range *cs {
		if a.GetName() != k {
			continue
		}

		digest = v
		break
	}

	if digest == "" {
		return "", fmt.Errorf("could not find digest for asset %s in checksums", a.GetName())
	}

	if a.GetDigest() == "" {
		return "", fmt.Errorf("asset %s has no digest set, cannot verify", a.GetName())
	}

	if !strings.HasSuffix(a.GetDigest(), digest) {
		return "", fmt.Errorf("asset digest mismatch: expected %s, got %s", digest, a.GetDigest())
	}

	return digest, nil
}

func downloadFile(ctx context.Context, a *github.ReleaseAsset) (io.ReadCloser, error) {
	fmt.Println("⬇️ Downloading file", a.GetName())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.GetBrowserDownloadURL(), nil)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get: invalid status code %d", resp.StatusCode)
	}
	return resp.Body, nil
}

func update(ctx context.Context, v string, a *github.ReleaseAsset, digest string) error {
	resp, err := downloadFile(ctx, a)
	if err != nil {
		return err
	}

	h, err := hex.DecodeString(digest)
	if err != nil {
		return fmt.Errorf("invalid digest: %w", err)
	}

	defer resp.Close() //nolint
	fmt.Println("📦 Updating binary to", v)
	err = selfupdate.Apply(resp, selfupdate.Options{
		Checksum: h,
	})
	if err != nil {
		return fmt.Errorf("apply update: %w", err)
	}
	fmt.Println("✅ Done")
	return nil
}

func downloadBundle(ctx context.Context, a *github.ReleaseAsset) (*bundle.Bundle, error) {
	resp, err := downloadFile(ctx, a)
	if err != nil {
		return nil, err
	}
	defer resp.Close() //nolint
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp); err != nil {
		return nil, fmt.Errorf("could not read bundle: %w", err)
	}

	var bundle bundle.Bundle
	if err := bundle.UnmarshalJSON(buf.Bytes()); err != nil {
		return nil, fmt.Errorf("could not unmarshal bundle: %w", err)
	}

	return &bundle, nil
}

func downloadAndVerifyChecksum(ctx context.Context, a *github.ReleaseAsset, b *bundle.Bundle) (*map[string]string, error) {
	resp, err := downloadFile(ctx, a)
	if err != nil {
		return nil, err
	}
	defer resp.Close() //nolint

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp); err != nil {
		return nil, fmt.Errorf("could not read checksums: %w", err)
	}
	if err := checkBundle(ctx, buf, b); err != nil {
		return nil, fmt.Errorf("could not verify singuature: %w", err)
	}
	checksums := make(map[string]string)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, "  ")
		if len(split) != 2 {
			return nil, fmt.Errorf("invalid checksum line: %s", line)
		}

		checksums[split[1]] = split[0]
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("could not scan checksums: %w", err)
	}

	return &checksums, nil
}

func checkBundle(ctx context.Context, a *bytes.Buffer, b *bundle.Bundle) error {
	opts := tuf.DefaultOptions()
	client, err := tuf.New(opts)
	if err != nil {
		return err
	}

	trustedMaterial, err := root.GetTrustedRoot(client)
	if err != nil {
		return err
	}

	sev, err := verify.NewVerifier(trustedMaterial,
		verify.WithSignedCertificateTimestamps(1),
		verify.WithTransparencyLog(1),
		verify.WithObserverTimestamps(1))
	if err != nil {
		return err
	}

	certID, err := verify.NewShortCertificateIdentity(
		"https://token.actions.githubusercontent.com",
		"",
		"",
		"^https://github.com/outscale/octl/",
	)
	if err != nil {
		return err
	}

	_, err = sev.Verify(b, verify.NewPolicy(verify.WithArtifact(a), verify.WithCertificateIdentity(certID)))
	if err != nil {
		return err
	}

	fmt.Println("🔐 Signature verified successfully")
	return nil
}
