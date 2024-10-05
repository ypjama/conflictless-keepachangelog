package conflictless

import "fmt"

// FlagCollection is a collection of flags.
type FlagCollection struct {
	Bump             *string
	ChangelogFile    *string
	Command          string
	Directory        *string
	SkipVersionLinks bool
}

// Config is the configuration for the CLI.
type Config struct {
	Flags                FlagCollection
	Bump                 Bump
	ChangelogFile        string
	RepositoryConfigFile string
	RepositoryHeadFile   string
	Directory            string
	Changelog            *Changelog
	CreateExtension      string
}

func (cfg *Config) SetGenerateConfigsFromFlags() error {
	cfg.SetChangelogFileFromFlags()
	cfg.SetDirectoryFromFlags()

	return cfg.SetBumpFromFlags()
}

func (cfg *Config) SetCreateConfigsFromFlags() {
	cfg.SetDirectoryFromFlags()
}

func (cfg *Config) SetCheckConfigsFromFlags() {
	cfg.SetDirectoryFromFlags()
}

func (cfg *Config) SetDirectoryFromFlags() {
	if cfg.Flags.Directory != nil {
		cfg.Directory = *cfg.Flags.Directory
	}
}

func (cfg *Config) SetChangelogFileFromFlags() {
	if cfg.Flags.ChangelogFile != nil {
		cfg.ChangelogFile = *cfg.Flags.ChangelogFile
	}
}

// SetBumpFromFlags sets the bump type by parsing the flag string.
func (cfg *Config) SetBumpFromFlags() error {
	bumpFlag := ""

	if cfg.Flags.Bump != nil {
		bumpFlag = *cfg.Flags.Bump
	}

	switch bumpFlag {
	case "patch":
		cfg.Bump = BumpPatch
	case "minor":
		cfg.Bump = BumpMinor
	case "major":
		cfg.Bump = BumpMajor
	default:
		return fmt.Errorf("%w: %s", ErrInvalidBumpFlag, bumpFlag)
	}

	return nil
}
