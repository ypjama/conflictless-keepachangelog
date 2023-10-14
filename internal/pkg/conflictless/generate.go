package conflictless

import (
	"bytes"
	"conflictless-keepachangelog/pkg/schema"
	"fmt"
	"os"
)

func generate(cfg *config) {
	switch *cfg.flags.bump {
	case "patch":
		cfg.bump = bumpPatch
	case "minor":
		cfg.bump = bumpMinor
	case "major":
		cfg.bump = bumpMajor
	default:
		printErrorAndExit(fmt.Sprintf("invalid bump flag: %s", *cfg.flags.bump), usageGenerate)
	}

	changelogData, err := os.ReadFile(cfg.changelogFile)
	if err != nil {
		printErrorAndExit(err.Error(), usageGenerate)
	}

	combined, err := scanDir(*cfg.flags.directory)
	if err != nil {
		printErrorAndExit(err.Error(), usageGenerate)
	}

	fmt.Println(dataToMarkdown(&changelogData, combined))

	printUsageAndExit(cfg)
}

func dataToMarkdown(changelogData *[]byte, data *schema.Data) string {
	eol := getEndOfLineSequence(changelogData)

	out := ""

	out += changeSectionToMarkdown("Added", data.Added, eol)
	out += changeSectionToMarkdown("Changed", data.Changed, eol)
	out += changeSectionToMarkdown("Deprecated", data.Deprecated, eol)
	out += changeSectionToMarkdown("Removed", data.Removed, eol)
	out += changeSectionToMarkdown("Fixed", data.Fixed, eol)
	out += changeSectionToMarkdown("Security", data.Security, eol)

	return out
}

func changeSectionToMarkdown(section string, entries []string, eol string) string {
	if len(entries) == 0 {
		return ""
	}

	out := fmt.Sprintf("### %s%s%s", section, eol, eol)

	for _, entry := range entries {
		out += fmt.Sprintf("- %s%s", entry, eol)
	}

	out += eol

	return out
}

func getEndOfLineSequence(changelogData *[]byte) string {
	if !bytes.Contains(*changelogData, []byte{'\r', '\n'}) {
		return "\n"
	}

	return "\r\n"
}
