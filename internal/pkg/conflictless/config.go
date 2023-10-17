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

// SetBumpFromFlags sets the bump type by parsing the flag string.
func (cfg *Config) SetBumpFromFlags() error {
	switch *cfg.Flags.Bump {
	case "patch":
		cfg.Bump = BumpPatch
	case "minor":
		cfg.Bump = BumpMinor
	case "major":
		cfg.Bump = BumpMajor
	default:
		return fmt.Errorf("%w: %s", ErrInvalidBumpFlag, *cfg.Flags.Bump)
	}

	return nil
}
