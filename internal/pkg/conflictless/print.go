package conflictless

import (
	"fmt"
	"os"
)

func printUsageAndExit(cfg *config) {
	if cfg.flags.command == "" {
		printErrorAndExit("", usage)
	}

	switch cfg.flags.command {
	case commandCheck:
		usageCheck()
	case commandGen:
		usageGenerate()
	default:
		usage()
	}

	os.Exit(exitCodeMisuseError)
}

func printErrorAndExit(msg string, usageFunc func()) {
	if msg != "" {
		//nolint:forbidigo
		fmt.Printf("Error: %s\n\n", msg)
	}

	usageFunc()
	os.Exit(exitCodeMisuseError)
}
