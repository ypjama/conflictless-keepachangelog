package conflictless

// Create creates a new change-file.
func Create(cfg *Config) {
	err := cfg.SetCreateConfigsFromFlags()
	if err != nil {
		PrintErrorAndExit(err.Error(), usageCreateOnError)
	}

	filename, err := ParseCurrentGitBranchAsFilename(cfg)
	if err != nil {
		PrintErrorAndExit(err.Error(), usageCreateOnError)
	}

	cfg.ChangesFile = filename

	err = createChangeFile(cfg)
	if err != nil {
		PrintErrorAndExit(err.Error(), usageCreateOnError)
	}
}
