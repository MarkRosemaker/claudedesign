// Command claudedesign  extracts .
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MarkRosemaker/claudedesign"
	"github.com/spf13/afero"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}
}

func run() error {
	path := ""
	flag.StringVar(&path, "path", "", "path of the zipped design export")
	flag.Parse()

	if len(flag.Args()) != 1 {
		return fmt.Errorf("expects exactly one argument, got %v", flag.Args())
	}

	if path == "" {
		return fmt.Errorf("-path is required")
	}

	if err := claudedesign.Extract(path,
		afero.NewBasePathFs(afero.NewOsFs(), "frontend")); err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}

	return nil
}
