package conflictless

import (
	"fmt"
	"path/filepath"
)

// Create ... TODO: Write description.
func Create(cfg *Config) {
	cfg.SetCreateConfigsFromFlags()

	filename, err := ParseCurrentGitBranchAsFilename(cfg)
	if err != nil {
		PrintErrorAndExit(err.Error(), usageCreateOnError)
	}

	//nolint:godox
	// TODO: finish this function before submitting a merge request.

	//nolint:forbidigo
	fmt.Println(cfg.Directory)

	//nolint:forbidigo
	fmt.Println(filepath.Join(cfg.Directory, filename))
}
