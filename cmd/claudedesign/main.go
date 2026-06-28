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

	cmd := "unzip"
	switch len(flag.Args()) {
	case 0: //ok
	case 1:
		cmd = flag.Args()[0]
	default:
		return fmt.Errorf("expects at most one argument, got %v", flag.Args())
	}

	switch cmd {
	case "unzip", "extract":
		if path == "" {
			return fmt.Errorf("-path is required")
		}

		if err := claudedesign.Extract(path,
			afero.NewBasePathFs(afero.NewOsFs(), "frontend")); err != nil {
			fmt.Fprintf(os.Stderr, "warning: %v\n", err)
		}
	default:
		return fmt.Errorf("unknown command %q", cmd)
	}

	return nil
}
