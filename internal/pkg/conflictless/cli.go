package conflictless

import (
	"flag"
	"fmt"
	"os"
)

const (
	argIdxCommand        = 1
	argIdxHelpTopic      = 2
	exitCodeSuccess      = 0
	exitCodeGeneralError = 1
	exitCodeMisuseError  = 2
	commandCheck         = "check"
	commandGen           = "generate"
	commandHelp          = "help"
	defaultBump          = bumpMinor
	minArguments         = 2
)

func CLI() {
	cfg := config{
		flags: flagCollection{
			bump:      new(string),
			command:   "",
			directory: new(string),
		},
		bump:          defaultBump,
		changelogFile: "CHANGELOG.md",
	}
	parseCLIFlags(&cfg)

	if cfg.flags.command == "" {
		printUsageAndExit(&cfg)
	}

	switch cfg.flags.command {
	case commandCheck:
		printUsageAndExit(&cfg)
	case commandGen:
		generate(&cfg)
	case commandHelp:
		help()
	default:
		printErrorAndExit(fmt.Sprintf("invalid command: '%s'", cfg.flags.command), usage)
	}
}

func parseCLIFlags(cfg *config) {
	flag.Usage = usage

	if len(os.Args) > argIdxCommand {
		cfg.flags.command = os.Args[argIdxCommand]
	}

	var cmd *flag.FlagSet

	switch cfg.flags.command {
	case commandHelp:
		cmd = flag.NewFlagSet(commandHelp, flag.ExitOnError)
		cmd.Usage = usage
	case commandGen:
		cmd = flag.NewFlagSet(commandGen, flag.ExitOnError)
		cmd.Usage = usageGenerate

		defineBumpFlags(cfg, cmd)
		defineDirFlags(cfg, cmd)
	case commandCheck:
		cmd = flag.NewFlagSet(commandCheck, flag.ExitOnError)
		cmd.Usage = usageCheck

		defineDirFlags(cfg, cmd)
	}

	if cmd != nil {
		err := cmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
	}
}

func defineBumpFlags(cfg *config, fs *flag.FlagSet) {
	defaultBumpStr := "minor"

	fs.StringVar(cfg.flags.bump, "bump", defaultBumpStr, "")
	fs.StringVar(cfg.flags.bump, "b", defaultBumpStr, "")
}

func defineDirFlags(cfg *config, fs *flag.FlagSet) {
	defaultDir := "changes"

	fs.StringVar(cfg.flags.directory, "dir", defaultDir, "")
	fs.StringVar(cfg.flags.directory, "d", defaultDir, "")
}
