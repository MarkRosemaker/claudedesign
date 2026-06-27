// Package claudedesign extracts a Claude Design zip export into a directory
// suitable for embedding in a Go binary. It handles the common transformations:
// renaming the project HTML to index.html, routing specific files elsewhere,
// and cleaning stale files from the target directory.
package claudedesign

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const frontendDir = "frontend"

var ErrExportNotFound = fmt.Errorf("export not found, likely outdated export ID")

func Extract(exportID string) error {
	url := fmt.Sprintf("https://api.anthropic.com/v1/design/h/%s", exportID)

	fmt.Printf("Fetching %v...\n", url)
	rsp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer rsp.Body.Close() //nolint:errcheck

	switch rsp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return ErrExportNotFound
	default:
		return fmt.Errorf("got status %s", rsp.Status)
	}

	var extract func(io.Reader) error
	switch ct := rsp.Header.Get("Content-Type"); ct {
	case "application/gzip":
		extract = extractGzip
	case "application/tar":
		extract = extractTar
	default:
		return fmt.Errorf("unknown Content-Type: %q", ct)
	}

	// Clean out the directory before unzipping
	if err := os.RemoveAll(frontendDir); err != nil {
		return fmt.Errorf("cleaning directory: %w", err)
	}

	if err := os.MkdirAll(frontendDir, 0o755); err != nil {
		return fmt.Errorf("creating new frontend directory")
	}

	return extract(rsp.Body)
}

func extractGzip(r io.Reader) error {
	gzReader, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	return extractTar(gzReader)
}

func extractTar(r io.Reader) error {
	tr := tar.NewReader(r)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if err := extractTarFile(tr, h); err != nil {
			return fmt.Errorf("extracting %q: %w", h.Name, err)
		}
	}
}

func extractTarFile(tr *tar.Reader, h *tar.Header) error {
	dest := filepath.Join(frontendDir, h.Name)

	switch h.Typeflag {
	case tar.TypeDir:
		if err := os.MkdirAll(dest, h.FileInfo().Mode()); err != nil {
			return err
		}
	case tar.TypeReg:
		if err := os.MkdirAll(filepath.Dir(dest), 0777); err != nil {
			return err
		}

		f, err := os.Create(dest)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(f, tr); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown type %q", h.Typeflag)
	}

	return nil
}
