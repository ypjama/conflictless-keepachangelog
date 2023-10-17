package conflictless

import (
	"fmt"
	"os"
)

func help() {
	var topic string
	if len(os.Args) > argIdxHelpTopic {
		topic = os.Args[argIdxHelpTopic]
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
		os.Exit(exitCodeMisuseError)
	}

	os.Exit(exitCodeSuccess)
}
