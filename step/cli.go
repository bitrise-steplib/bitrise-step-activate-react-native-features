package step

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	cliVersion    = "1.7.0"
	cliBinaryName = "bitrise-build-cache"
	cliInstallDir = "/usr/local/bin"

	downloadURLTemplate = "https://github.com/bitrise-io/bitrise-build-cache-cli/releases/download/v%s/bitrise-build-cache_%s_%s_%s.tar.gz"

	maxBinarySize = 500 << 20 // 500 MB safety limit
)

// installCLI downloads and installs the CLI binary from GitHub releases.
// If the correct version is already installed, it returns the existing path.
func installCLI(ctx context.Context, logger Logger) (string, error) {
	binaryPath := filepath.Join(cliInstallDir, cliBinaryName)

	if isCorrectVersion(binaryPath) {
		logger.Infof("CLI %s already installed at %s", cliVersion, binaryPath)

		return binaryPath, nil
	}

	url := fmt.Sprintf(downloadURLTemplate, cliVersion, cliVersion, runtime.GOOS, runtime.GOARCH)
	logger.Infof("Downloading CLI v%s from %s", cliVersion, url)

	if err := downloadAndInstall(ctx, url, binaryPath); err != nil {
		return "", fmt.Errorf("install CLI v%s: %w", cliVersion, err)
	}

	logger.Infof("CLI v%s installed at %s", cliVersion, binaryPath)

	return binaryPath, nil
}

func isCorrectVersion(binaryPath string) bool {
	out, err := exec.Command(binaryPath, "--version").Output() //nolint:gosec
	if err != nil {
		return false
	}

	return strings.Contains(string(out), cliVersion)
}

func downloadAndInstall(ctx context.Context, url, destPath string) error {
	body, err := downloadFile(ctx, url)
	if err != nil {
		return fmt.Errorf("download %s: %w", url, err)
	}
	defer body.Close()

	if err := extractBinaryFromTarGz(body, destPath); err != nil {
		return fmt.Errorf("extract binary: %w", err)
	}

	if err := os.Chmod(destPath, 0o755); err != nil {
		return fmt.Errorf("chmod: %w", err)
	}

	return nil
}

func downloadFile(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()

		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func extractBinaryFromTarGz(r io.Reader, destPath string) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("gzip reader: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			return fmt.Errorf("binary %s not found in archive", cliBinaryName)
		}

		if err != nil {
			return fmt.Errorf("read tar entry: %w", err)
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		if filepath.Base(header.Name) != cliBinaryName {
			continue
		}

		f, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
		if err != nil {
			return fmt.Errorf("create file: %w", err)
		}
		defer f.Close()

		if _, err := io.Copy(f, io.LimitReader(tr, maxBinarySize)); err != nil {
			return fmt.Errorf("write binary: %w", err)
		}

		return nil
	}
}
