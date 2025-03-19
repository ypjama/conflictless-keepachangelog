package conflictless

import (
	"fmt"
	"strings"
)

// FlagCollection is a collection of flags.
type FlagCollection struct {
	Bump                     *string
	ChangeFileFormat         *string
	ChangeFileName           *string
	ChangelogFile            *string
	ChangeTypesCsv           *string
	Command                  string
	Directory                *string
	SkipVersionLinks         bool
	UseVPrefixInVersionLinks bool
}

// Config is the configuration for the CLI.
type Config struct {
	Bump                 Bump
	ChangeFileFormat     string
	Changelog            *Changelog
	ChangelogFile        string
	ChangeFile           string
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
	cfg.SetChangeFileFromFlags()

	return cfg.SetChangeFileFormatFromFlags()
}

func (cfg *Config) SetCheckConfigsFromFlags() {
	cfg.SetDirectoryFromFlags()
}

func (cfg *Config) SetPreviewConfigsFromFlags() error {
	cfg.SetChangelogFileFromFlags()
	cfg.SetDirectoryFromFlags()

	return cfg.SetBumpFromFlags()
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

func (cfg *Config) SetChangeFileFromFlags() {
	if cfg.Flags.ChangeFileName != nil {
		cfg.ChangeFile = *cfg.Flags.ChangeFileName
	}
}

func (cfg *Config) SetChangeFileFormatFromFlags() error {
	if cfg.Flags.ChangeFileFormat == nil {
		return nil
	}

	formatFlag := *cfg.Flags.ChangeFileFormat
	formatFlag = strings.ToLower(formatFlag)

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

func (cfg *Config) GenerateNewSection() string {
	var err error

	cfg.Changelog, err = ReadChangelog(cfg)
	if err != nil {
		PrintErrorAndExit(err.Error(), func() {})
	}

	combined, err := scanDir(cfg.Directory)
	if err != nil {
		PrintErrorAndExit(err.Error(), func() {})
	}

	if combined.IsEmpty() {
		PrintErrorAndExit("no changelog entries found", func() {})
	}

	newSection := DataToMarkdown(cfg, combined)
	if newSection == "" {
		PrintErrorAndExit("failed to generate a new version section", func() {})
	}

	return newSection
}
