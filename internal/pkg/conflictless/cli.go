package conflictless

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	argIdxCommand           = 1
	argIdxHelpTopic         = 2
	commandCheck            = "check"
	commandCreate           = "create"
	commandGen              = "generate"
	commandHelp             = "help"
	defaultBump             = BumpMinor
	defaultChangeFileFormat = "yml"
	defaultChangeTypesCSV   = "changed"
	defaultDirectory        = "changes"
	exitCodeGeneralError    = 1
	exitCodeMisuseError     = 2
	exitCodeSuccess         = 0
	minArguments            = 2
)

func CLI() {
	cfg := Config{
		Flags: FlagCollection{
			Bump:             new(string),
			ChangeFileFormat: new(string),
			ChangelogFile:    new(string),
			ChangeTypesCsv:   new(string),
			Command:          "",
			Directory:        new(string),
			SkipVersionLinks: false,
		},
		Bump:                 defaultBump,
		Changelog:            nil,
		ChangelogFile:        "CHANGELOG.md",
		ChangesFile:          "",
		ChangeTypesCsv:       defaultChangeTypesCSV,
		ChangeFileFormat:     defaultChangeFileFormat,
		Directory:            defaultDirectory,
		RepositoryConfigFile: ".git/config",
		RepositoryHeadFile:   ".git/HEAD",
	}
	parseCLIFlags(&cfg)

	if cfg.Flags.Command == "" {
		PrintUsageAndExit(&cfg)
	}

	switch cfg.Flags.Command {
	case commandCheck:
		Check(&cfg)
	case commandCreate:
		Create(&cfg)
	case commandGen:
		Generate(&cfg)
	case commandHelp:
		Help()
	default:
		PrintErrorAndExit(fmt.Sprintf("invalid command: '%s'", cfg.Flags.Command), usageOnError)
	}
}

func argsWithoutTestFlags() []string {
	args := []string{}

	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			continue
		}

		args = append(args, arg)
	}

	return args
}

func parseCLIFlags(cfg *Config) {
	flag.Usage = usage

	args := argsWithoutTestFlags()

	if len(args) > argIdxCommand {
		cfg.Flags.Command = args[argIdxCommand]
	}

	var cmd *flag.FlagSet

	switch cfg.Flags.Command {
	case commandHelp:
		cmd = flag.NewFlagSet(commandHelp, flag.ExitOnError)
		cmd.Usage = usageOnError
	case commandGen:
		cmd = flag.NewFlagSet(commandGen, flag.ExitOnError)
		cmd.Usage = usageGenerateOnError

		defineBumpFlags(cfg, cmd)
		defineChangeLogFlags(cfg, cmd)
		defineDirFlags(cfg, cmd)
		defineSkipFlags(cfg, cmd)
	case commandCheck:
		cmd = flag.NewFlagSet(commandCheck, flag.ExitOnError)
		cmd.Usage = usageCheckOnError

		defineDirFlags(cfg, cmd)
	case commandCreate:
		cmd = flag.NewFlagSet(commandCreate, flag.ExitOnError)
		cmd.Usage = usageCreateOnError

		defineDirFlags(cfg, cmd)
	}

	if cmd != nil {
		err := cmd.Parse(args[2:])
		if err != nil {
			panic(err)
		}
	}
}

func defineChangeLogFlags(cfg *Config, fs *flag.FlagSet) {
	defaultChangelogFile := "CHANGELOG.md"

	fs.StringVar(cfg.Flags.ChangelogFile, "changelog", defaultChangelogFile, "")
	fs.StringVar(cfg.Flags.ChangelogFile, "c", defaultChangelogFile, "")
}

func defineBumpFlags(cfg *Config, fs *flag.FlagSet) {
	defaultBumpStr := "minor"

	fs.StringVar(cfg.Flags.Bump, "bump", defaultBumpStr, "")
	fs.StringVar(cfg.Flags.Bump, "b", defaultBumpStr, "")
}

func defineDirFlags(cfg *Config, fs *flag.FlagSet) {
	defaultDir := "changes"

	fs.StringVar(cfg.Flags.Directory, "dir", defaultDir, "")
	fs.StringVar(cfg.Flags.Directory, "d", defaultDir, "")
}

func defineSkipFlags(cfg *Config, fs *flag.FlagSet) {
	fs.BoolVar(&cfg.Flags.SkipVersionLinks, "skip-version-links", false, "")
	fs.BoolVar(&cfg.Flags.SkipVersionLinks, "s", false, "")
}
