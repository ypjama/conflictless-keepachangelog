package conflictless_test

import (
	"conflictless-keepachangelog/internal/pkg/conflictless"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const changelogContent = `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.1] - 2023-06-26

### Security

- whoops

## [1.0.0] - 2023-06-24

### Changed

- a lot

## [0.2.0] - 2023-03-05

### Added

- foo
- bar

### Changed

- baz

## [0.1.0] - 2022-11-08

### Added

- foo

### Changed

- bar
- baz
`

func TestNextReleaseHeader(t *testing.T) {
	t.Parallel()

	file, err := os.CreateTemp(os.TempDir(), "CHANGELOG.md")
	assert.NoError(t, err)

	defer os.Remove(file.Name())

	_, err = file.WriteString(changelogContent)
	assert.NoError(t, err)

	cfg := new(conflictless.Config)
	cfg.ChangelogFile = file.Name()

	changelogContents, err := conflictless.ReadChangelog(cfg)

	assert.NoError(t, err)
	assert.Equal(t, "1.0.2", changelogContents.NextReleaseHeader(conflictless.BumpPatch))
	assert.Equal(t, "1.1.0", changelogContents.NextReleaseHeader(conflictless.BumpMinor))
	assert.Equal(t, "2.0.0", changelogContents.NextReleaseHeader(conflictless.BumpMajor))
}
