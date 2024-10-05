package conflictless

import "fmt"

// FlagCollection is a collection of flags.
type FlagCollection struct {
	Bump             *string
	ChangeFileFormat *string
	ChangelogFile    *string
	ChangeTypesCsv   *string
	Command          string
	Directory        *string
	SkipVersionLinks bool
}

// Config is the configuration for the CLI.
type Config struct {
	Bump                 Bump
	ChangeFileFormat     string
	Changelog            *Changelog
	ChangelogFile        string
	ChangesFile          string
	ChangeTypesCsv       string
	Directory            string
	Flags                FlagCollection
	RepositoryConfigFile string
	RepositoryHeadFile   string
}

func (cfg *Config) SetGenerateConfigsFromFlags() error {
	cfg.SetChangelogFileFromFlags()
	cfg.SetDirectoryFromFlags()

	return cfg.SetBumpFromFlags()
}

func (cfg *Config) SetCreateConfigsFromFlags() error {
	cfg.SetChangeTypesFromFlags()
	cfg.SetDirectoryFromFlags()

	return cfg.SetChangeFileFormatFromFlags()
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

func (cfg *Config) SetChangeTypesFromFlags() {
	if cfg.Flags.ChangeTypesCsv != nil {
		cfg.ChangeTypesCsv = *cfg.Flags.ChangeTypesCsv
	}
}

func (cfg *Config) SetChangeFileFormatFromFlags() error {
	if cfg.Flags.ChangeFileFormat == nil {
		return nil
	}

	formatFlag := *cfg.Flags.ChangeFileFormat

	switch formatFlag {
	case "yaml":
		cfg.ChangeFileFormat = "yaml"
	case "yml":
		cfg.ChangeFileFormat = "yml"
	case "json":
		cfg.ChangeFileFormat = "json"
	default:
		return fmt.Errorf("%w, %s", ErrInvalidFormatFlag, formatFlag)
	}

	return nil
}
