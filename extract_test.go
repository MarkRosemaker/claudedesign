package claudedesign

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"testing"
)

//go:embed testdata/Hg869Yo0PbGgmpUWg3hNdw.tar.gz
var data []byte

type testTransport struct {
}

// RoundTrip implements [http.RoundTripper].
func (t *testTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.URL.String() != "https://api.anthropic.com/v1/design/h/Hg869Yo0PbGgmpUWg3hNdw" {
		return nil, fmt.Errorf("unknown request url %q", r.URL)
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(data)),
		Header: http.Header{
			"Content-Type": []string{"application/gzip"},
		},
	}, nil
}

func TestExtract(t *testing.T) {
	http.DefaultTransport = &testTransport{}

	if err := Extract("Hg869Yo0PbGgmpUWg3hNdw"); err != nil {
		t.Fatal(err)
	}
	// Hg869Yo0PbGgmpUWg3hNdw
	// Simulate a Claude Design export: MyApp.zip containing MyApp.html + app.jsx
	// makeZip(t, filepath.Join(dir, "MyApp.zip"), map[string]string{
	// 	"MyApp.html": "<html>hello</html>",
	// 	"app.jsx":    "const App = () => {};",
	// })

	// // Place a stale file and a favicon that should be preserved.
	// writeFile(t, filepath.Join(dir, "old-file.js"), []byte("stale"))
	// writeFile(t, filepath.Join(dir, "favicon.ico"), []byte("icon"))

	// if err := Extract(dir); err != nil {
	// 	t.Fatal(err)
	// }

	// // index.html should exist (renamed from MyApp.html)
	// assertFile(t, filepath.Join(dir, "index.html"), "<html>hello</html>")

	// // app.jsx extracted as-is
	// assertFile(t, filepath.Join(dir, "app.jsx"), "const App = () => {};")

	// // favicon preserved
	// assertFile(t, filepath.Join(dir, "favicon.ico"), "icon")

	// // stale file removed
	// if _, err := os.Stat(filepath.Join(dir, "old-file.js")); err == nil {
	// 	t.Error("old-file.js should have been deleted")
	// }
}
