package claudedesign

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"errors"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

//go:embed "testdata/LSO Helper.zip"
var data []byte

func TestExtract(t *testing.T) {
	mem := &afero.MemMapFs{}

	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}

	if err := extract(r, mem); err != nil {
		t.Fatal(err)
	}

	assertFile(t, mem, "index.html", "<title>LSO Helper</title>")
	assertFile(t, mem, "openapi.json", `"openapi": "3.0.0",`)
}

func assertFile(t *testing.T, mem *afero.MemMapFs, path, containing string) {
	t.Helper()

	got, err := afero.ReadFile(mem, path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			if entries, err2 := afero.ReadDir(mem, filepath.Dir(path)); err2 == nil {
				names := make([]string, len(entries))
				for i, e := range entries {
					names[i] = e.Name()
				}

				t.Fatalf("%v but %q exist", err, names)
			}
		}

		t.Fatal(err)
	}

	if !bytes.Contains(got, []byte(containing)) {
		t.Errorf("content=%s, want=%s", got, containing)
	}
}
