package conflictless

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	minSemverParts = 3
)

type Changelog struct {
	Filepath       string
	ReleaseHeaders []string
	Bytes          []byte
}

func ReadChangelog(cfg *Config) (*Changelog, error) {
	changelog := new(Changelog)

	data, err := os.ReadFile(cfg.ChangelogFile)
	if err != nil {
		return nil, fmt.Errorf("%w. %w", ErrChangelogFileNotFound, err)
	}

	changelog.Filepath = cfg.ChangelogFile
	changelog.Bytes = data
	changelog.ReleaseHeaders = ParseReleaseHeaders(data)

	return changelog, nil
}

// LatestReleaseHeader returns the latest release header.
func (cc *Changelog) LatestReleaseHeader() string {
	latest := ""
	latestInt := 0

	for _, header := range cc.ReleaseHeaders {
		major, minor, patch := semverTagToIntegers(header)

		semverInt := major*1000000 + minor*1000 + patch
		if semverInt > latestInt {
			latest = header
			latestInt = semverInt
		}
	}

	return latest
}

// NextReleaseHeader returns the next release header based on the latest release header and the bump type.
func (cc *Changelog) NextReleaseHeader(bump Bump) string {
	latest := cc.LatestReleaseHeader()
	if latest == "" || strings.ToLower(latest) == "unreleased" {
		return bump.InitialVersion()
	}

	major, minor, patch := semverTagToIntegers(latest)
	if major == 0 && minor == 0 && patch == 0 {
		return bump.InitialVersion()
	}

	switch bump {
	case BumpPatch:
		patch++
	case BumpMinor:
		minor++

		patch = 0
	case BumpMajor:
		major++

		minor = 0
		patch = 0
	}

	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

func (cc *Changelog) WriteSection(section string) error {
	latest := cc.LatestReleaseHeader()

	if latest == "" {
		beforeIdx := len(string(cc.Bytes))
		re := regexp.MustCompile(`(?i)##\s*\[unreleased]\s+(?:.*\n)*?##\s*\[([\d.]+)\]`)

		idx := re.FindSubmatchIndex(cc.Bytes)
		if len(idx) > 1 {
			beforeIdx = idx[0]
		}

		return writeBeforeIndex(cc, beforeIdx, section, true)
	}

	re, err := regexp.Compile(`##\s*\[` + latest + `\]`)
	if err != nil {
		return fmt.Errorf("%w. %w", ErrChangelogWrite, err)
	}

	idx := re.FindIndex(cc.Bytes)
	if len(idx) < 1 {
		return fmt.Errorf("%w - could not find section %s", ErrChangelogWrite, latest)
	}

	return writeBeforeIndex(cc, idx[0], section, false)
}

func writeBeforeIndex(changelog *Changelog, beforeIdx int, section string, startWithNewline bool) error {
	var updatedContent []byte

	if startWithNewline {
		section = getEndOfLineSequence(&changelog.Bytes) + section
	}

	content := string(changelog.Bytes)

	if beforeIdx < 1 {
		updatedContent = []byte(section + content)
	} else {
		updatedContent = []byte(content[:beforeIdx] + section + content[beforeIdx:])
	}

	err := os.WriteFile(changelog.Filepath, updatedContent, fs.FileMode(writeFileMode))
	if err != nil {
		return fmt.Errorf("%w. %w", ErrChangelogWrite, err)
	}

	return nil
}

func semverTagToIntegers(semver string) (int, int, int) {
	slices := strings.Split(semver, ".")
	if len(slices) < minSemverParts {
		return 0, 0, 0
	}

	major, err := strconv.Atoi(slices[0])
	if err != nil {
		return 0, 0, 0
	}

	minor, err := strconv.Atoi(slices[1])
	if err != nil {
		return 0, 0, 0
	}

	patch, err := strconv.Atoi(slices[2])
	if err != nil {
		return 0, 0, 0
	}

	return major, minor, patch
}
