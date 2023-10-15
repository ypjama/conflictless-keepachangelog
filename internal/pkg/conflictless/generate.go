package conflictless

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ypjama/conflictless-keepachangelog/pkg/schema"
)

func generate(cfg *Config) {
	switch *cfg.Flags.Bump {
	case "patch":
		cfg.Bump = BumpPatch
	case "minor":
		cfg.Bump = BumpMinor
	case "major":
		cfg.Bump = BumpMajor
	default:
		printErrorAndExit(fmt.Sprintf("invalid bump flag: %s", *cfg.Flags.Bump), usageGenerate)
	}

	combined, err := scanDir(*cfg.Flags.Directory)
	if err != nil {
		printErrorAndExit(err.Error(), func() {})
	}

	if combined.IsEmpty() {
		printErrorAndExit("no changelog entries found", func() {})
	}

	newSection := DataToMarkdown(cfg, combined)
	if newSection == "" {
		printErrorAndExit("failed to generate a new version section", func() {})
	}

	err = cfg.Changelog.WriteSection(newSection)
	if err != nil {
		printErrorAndExit(err.Error(), func() {})
	}

	err = removeChangeFiles(*cfg.Flags.Directory)
	if err != nil {
		printErrorAndExit(err.Error(), func() {})
	}

	printGenerateSuccess(newSection)
}

func DataToMarkdown(cfg *Config, data *schema.Data) string {
	if cfg.Changelog == nil {
		return ""
	}

	eol := getEndOfLineSequence(&cfg.Changelog.Bytes)

	out := ""

	sectionName := cfg.Changelog.NextReleaseHeader(cfg.Bump)
	dateStr := time.Now().Format("2006-01-02")

	var sectionLink string

	if !cfg.Flags.SkipVersionLinks {
		sectionLink = SectionLink(ParseRepositoryURL(cfg), sectionName)
	}

	out += fmt.Sprintf("## ["+sectionName+"] - %s%s%s", dateStr, eol, eol)

	if sectionLink != "" {
		out += fmt.Sprintf("%s%s%s", sectionLink, eol, eol)
	}

	out += entriesToMarkdown("Added", data.Added, eol)
	out += entriesToMarkdown("Changed", data.Changed, eol)
	out += entriesToMarkdown("Deprecated", data.Deprecated, eol)
	out += entriesToMarkdown("Removed", data.Removed, eol)
	out += entriesToMarkdown("Fixed", data.Fixed, eol)
	out += entriesToMarkdown("Security", data.Security, eol)

	return out
}

func entriesToMarkdown(section string, entries []string, eol string) string {
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
