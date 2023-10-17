package conflictless

func check(cfg *Config) {
	combined, err := scanDir(*cfg.Flags.Directory)
	if err != nil {
		PrintErrorAndExit(err.Error(), func() {})
	}

	PrintCheckSuccess(combined.IsEmpty())
}
