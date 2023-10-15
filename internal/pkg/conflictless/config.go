package conflictless

// FlagCollection is a collection of flags.
type FlagCollection struct {
	Bump             *string
	Directory        *string
	Command          string
	SkipVersionLinks bool
}

// Config is the configuration for the CLI.
type Config struct {
	Flags                FlagCollection
	Bump                 Bump
	ChangelogFile        string
	RepositoryConfigFile string
	Changelog            *Changelog
}
