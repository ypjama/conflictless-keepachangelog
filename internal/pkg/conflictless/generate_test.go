package conflictless_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ypjama/conflictless-keepachangelog/internal/pkg/conflictless"
)

const (
	mkdirFileMode = 0o755
	writeFileMode = 0o644
)

func writeDataToFile(t *testing.T, data []byte, file *os.File) {
	t.Helper()

	err := os.WriteFile(file.Name(), data, writeFileMode)
	if err != nil {
		t.Fatal(err)
	}
}

func createFile(t *testing.T, dir, name string) *os.File {
	t.Helper()

	file, err := os.Create(dir + "/" + name)
	if err != nil {
		t.Fatal(err)
	}

	return file
}

func createTempFile(t *testing.T, dir, pattern string) *os.File {
	t.Helper()

	file, err := os.CreateTemp(dir, pattern)
	if err != nil {
		t.Fatal(err)
	}

	return file
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	changesDir, err := os.MkdirTemp(os.TempDir(), "changes")
	assert.NoError(t, err)

	defer os.RemoveAll(changesDir)

	os.TempDir()

	changesFile := createFile(t, changesDir, "1-foo.json")
	defer os.Remove(changesFile.Name())

	changelogFile := createTempFile(t, os.TempDir(), "CHANGELOG.md")
	defer os.Remove(changelogFile.Name())

	gitConfigFile := createTempFile(t, os.TempDir(), ".gitconfig")
	defer os.Remove(gitConfigFile.Name())

	writeDataToFile(t, []byte(`{"fixed":["foo"]}`), changesFile)
	writeDataToFile(t, []byte(changelogContent), changelogFile)
	writeDataToFile(t, []byte(gitConfig), gitConfigFile)

	flagValueBumpPatch := "patch"

	cfg := new(conflictless.Config)
	cfg.Flags.Directory = &changesDir
	cfg.Flags.Bump = &flagValueBumpPatch
	cfg.ChangelogFile = changelogFile.Name()
	cfg.RepositoryConfigFile = gitConfigFile.Name()

	conflictless.Generate(cfg)

	actual, err := os.ReadFile(changelogFile.Name())
	assert.NoError(t, err)
	assert.Contains(t, string(actual), "## [Unreleased]")
	assert.Contains(t, string(actual), "## [1.0.2]")
}
