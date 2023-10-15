package conflictless

import (
	"fmt"
	"os"
)

func printUsageAndExit(cfg *Config) {
	if cfg.Flags.Command == "" {
		printErrorAndExit("", usage)
	}

	switch cfg.Flags.Command {
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

func printGenerateSuccess(section string) {
	//nolint:forbidigo
	fmt.Printf("Generated new version section successfully!\n\n```\n%s```\n", section)
}
