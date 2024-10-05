package conflictless

// Create creates a new change-file.
func Create(cfg *Config) {
	err := cfg.SetCreateConfigsFromFlags()
	if err != nil {
		PrintErrorAndExit(err.Error(), usageCreateOnError)
	}

	if cfg.ChangeFile == "" {
		filename, err := ParseCurrentGitBranchAsFilename(cfg)
		if err != nil {
			PrintErrorAndExit(err.Error(), usageCreateOnError)
		}

		cfg.ChangeFile = filename
	} else {
		cfg.ChangeFile += "." + cfg.ChangeFileFormat
	}

	err = createChangeFile(cfg)
	if err != nil {
		PrintErrorAndExit(err.Error(), usageCreateOnError)
	}

	PrintCreateSuccess(cfg)
}
