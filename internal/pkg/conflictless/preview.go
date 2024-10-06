package conflictless

// Preview prints a  preview of unreleased changes.
func Preview(cfg *Config) {
	err := cfg.SetPreviewConfigsFromFlags()
	if err != nil {
		PrintErrorAndExit(err.Error(), usagePreviewOnError)
	}

	PrintPreviewSuccess(cfg.GenerateNewSection())
}
