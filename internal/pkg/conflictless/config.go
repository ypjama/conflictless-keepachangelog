package conflictless

import "fmt"

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

func (cfg *Config) setBumpFromFlags() error {
	switch *cfg.Flags.Bump {
	case "patch":
		cfg.Bump = BumpPatch
	case "minor":
		cfg.Bump = BumpMinor
	case "major":
		cfg.Bump = BumpMajor
	default:
		return fmt.Errorf("%w: %s", errInvalidBumpFlag, *cfg.Flags.Bump)
	}

	return nil
}
