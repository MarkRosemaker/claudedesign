// Command claudedesign  extracts .
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/MarkRosemaker/claudedesign"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}
}

func run() error {
	exportID := ""
	flag.StringVar(&exportID, "id", "", "id of the design: https://api.anthropic.com/v1/design/h/{id}")
	flag.Parse()

	if len(flag.Args()) != 1 {
		return fmt.Errorf("expects exactly one argument, got %v", flag.Args())
	}

	switch arg := flag.Args()[0]; arg {
	case "unzip", "extract", "import":
		if exportID == "" {
			return fmt.Errorf("-id is required")
		}

		if err := claudedesign.Extract(exportID); errors.Is(err, claudedesign.ErrExportNotFound) {
			fmt.Fprintf(os.Stderr, "warning: %v\n", err)
		} else if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown argument %q", arg)
	}

	return nil
}
