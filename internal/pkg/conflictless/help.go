package conflictless

import (
	"fmt"
	"os"
)

func help() {
	var topic string

	args := argsWithoutTestFlags()

	if len(args) > argIdxHelpTopic {
		topic = args[argIdxHelpTopic]
	}

	switch topic {
	case commandCheck:
		usageCheck()
	case commandGen:
		usageGenerate()
	case "":
		usage()
	default:
		PrintErrorAndExit(fmt.Sprintf("unknown help topic '%s'", topic), func() {})
	}

	os.Exit(exitCodeSuccess)
}
