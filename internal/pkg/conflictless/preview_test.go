package conflictless_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ypjama/conflictless-keepachangelog/internal/pkg/conflictless"
)

func TestPreview(t *testing.T) {
	t.Parallel()

	changesDir, err := os.MkdirTemp(os.TempDir(), "changes")
	assert.NoError(t, err)

	defer os.RemoveAll(changesDir)

	changesFile := createFile(t, changesDir, "test-preview.json")
	defer os.Remove(changesFile.Name())

	changelogFile := createTempFile(t, os.TempDir(), "test-preview-CHANGELOG.md")
	defer os.Remove(changelogFile.Name())

	gitConfigFile := createTempFile(t, os.TempDir(), "test-preview.gitconfig")
	defer os.Remove(gitConfigFile.Name())

	writeDataToFile(t, []byte(`{"added":["New major feature"]}`), changesFile)
	writeDataToFile(t, []byte(changelogContent), changelogFile)
	writeDataToFile(t, []byte(gitConfig), gitConfigFile)

	flagValueBumpPatch := "major"

	cfg := new(conflictless.Config)
	cfg.Flags.Directory = &changesDir
	cfg.Flags.Bump = &flagValueBumpPatch
	cfg.Flags.SkipVersionLinks = true
	cfg.ChangelogFile = changelogFile.Name()
	cfg.RepositoryConfigFile = gitConfigFile.Name()

	startStdoutCapture(t)

	conflictless.Preview(cfg)

	output := stopStdoutCapture(t)

	assert.Contains(t, output, "```md\n")
	assert.Contains(t, output, "## [2.0.0]")
	assert.Contains(t, output, "### Added\n\n- New major feature\n")
}
