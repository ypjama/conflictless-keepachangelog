package conflictless

func check(cfg *Config) {
	combined, err := scanDir(*cfg.Flags.Directory)
	if err != nil {
		printErrorAndExit(err.Error(), func() {})
	}

	printCheckSuccess(combined.IsEmpty())
}