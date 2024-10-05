package conflictless

import (
	"fmt"
	"os"
	"path/filepath"
)

// PrintUsageAndExit prints the usage and exits.
func PrintUsageAndExit(cfg *Config) {
	if cfg.Flags.Command == "" {
		PrintErrorAndExit("", usageOnError)
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

// PrintErrorAndExit prints an error message and exits.
func PrintErrorAndExit(msg string, usageFunc func()) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", msg)
	}

	usageFunc()
	os.Exit(exitCodeMisuseError)
}

// PrintGenerateSuccess prints a success message for generate command.
func PrintGenerateSuccess(section string) {
	//nolint:forbidigo
	fmt.Printf("Generated new version section successfully!\n\n```md\n%s```\n", section)
}

// PrintCheckSuccess prints a success message for check command.
func PrintCheckSuccess(noContent bool) {
	var msg string

	if noContent {
		msg = "No changes found!"
	} else {
		msg = "Change files are valid!"
	}

	//nolint:forbidigo
	fmt.Println(msg)
}

func PrintCreateSuccess(cfg *Config) {
	//nolint:forbidigo
	fmt.Printf("Created new change-file '%s' successfully!\n", filepath.Join(cfg.Directory, cfg.ChangeFile))
}
