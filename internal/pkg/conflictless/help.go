package conflictless

import (
	"fmt"
	"os"
)

// Help prints the help message and exits.
func Help() {
	var topic string

	args := argsWithoutTestFlags()

	if len(args) > argIdxHelpTopic {
		topic = args[argIdxHelpTopic]
	}

	switch topic {
	case commandCheck:
		usageCheck()
	case commandCreate:
		usageCreate()
	case commandGen:
		usageGenerate()
	case "":
		usage()
	default:
		PrintErrorAndExit(fmt.Sprintf("unknown help topic '%s'", topic), func() {})
	}

	os.Exit(exitCodeSuccess)
}
