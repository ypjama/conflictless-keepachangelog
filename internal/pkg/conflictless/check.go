package conflictless

// Check checks the validity of the change files.
func Check(cfg *Config) {
	combined, err := scanDir(*cfg.Flags.Directory)
	if err != nil {
		PrintErrorAndExit(err.Error(), func() {})
	}

	PrintCheckSuccess(combined.IsEmpty())
}
