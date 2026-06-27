// Package claudedesign extracts a Claude Design zip export into a directory
// suitable for embedding in a Go binary. It handles the common transformations:
// renaming the project HTML to index.html, routing specific files elsewhere,
// and cleaning stale files from the target directory.
package claudedesign

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

const frontendDir = "frontend"

func Extract(zipFile string, fsys afero.Fs) error {
	// Open the zip archive for reading
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer r.Close()

	// Clean out the directory before unzipping
	if err := fsys.RemoveAll(frontendDir); err != nil {
		return fmt.Errorf("cleaning directory: %w", err)
	}

	if err := fsys.MkdirAll(frontendDir, 0o755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	if err := extract(&r.Reader, afero.NewBasePathFs(fsys, frontendDir)); err != nil {
		return err
	}

	if err := os.Remove(zipFile); err != nil {
		return err
	}

	return nil
}

func extract(r *zip.Reader, fsys afero.Fs) error {

	// Iterate over each file/dir in the archive
	for _, f := range r.File {
		if err := extractFile(fsys, f); err != nil {
			return fmt.Errorf("extracting file %s: %w", f.Name, err)
		}
	}

	return nil
}

func extractFile(fsys afero.Fs, f *zip.File) error {
	if f.FileInfo().IsDir() {
		// Create directory
		if err := fsys.MkdirAll(f.Name, f.Mode()); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		return nil
	}

	// Ensure the parent directory exists (for files in subfolders)
	if err := fsys.MkdirAll(filepath.Dir(f.Name), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Open the file inside the archive
	src, err := f.Open()
	if err != nil {
		return fmt.Errorf("failed to open archived file: %w", err)
	}
	defer src.Close()

	// Create the destination file with the same permissions
	dst, err := fsys.OpenFile(f.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy the content
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}
