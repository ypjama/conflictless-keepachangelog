package conflictless

// Check checks the validity of the change files.
func Check(cfg *Config) {
	cfg.SetCheckConfigsFromFlags()

	combined, err := scanDir(cfg.Directory)
	if err != nil {
		PrintErrorAndExit(err.Error(), func() {})
	}

	PrintCheckSuccess(combined.IsEmpty())
}
