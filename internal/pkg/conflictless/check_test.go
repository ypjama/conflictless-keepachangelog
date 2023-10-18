package conflictless_test

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ypjama/conflictless-keepachangelog/internal/pkg/conflictless"
)

//nolint:paralleltest // this test is not parallel because it modifies os.Stdout.
func TestCheck(t *testing.T) {
	cfg := new(conflictless.Config)

	changesDir, err := os.MkdirTemp(os.TempDir(), "changes")
	assert.NoError(t, err)

	defer os.RemoveAll(changesDir)

	changesFile := createFile(t, changesDir, "test-check.json")
	defer os.Remove(changesFile.Name())

	stdOutFile := createTempFile(t, os.TempDir(), "stdout-test-check")
	defer os.Remove(stdOutFile.Name())

	writeDataToFile(t, []byte(`{"fixed":["foo"]}`), changesFile)

	cfg.Flags.Directory = &changesDir

	os.Stdout = stdOutFile

	conflictless.Check(cfg)

	data, err := os.ReadFile(stdOutFile.Name())

	assert.NoError(t, err)
	assert.Equal(t, "Change files are valid!\n", string(data))
}

func TestCheckWhenDirectoryDoesNotExist(t *testing.T) {
	t.Parallel()

	if os.Getenv("TEST_CHECK_WHEN_ERROR") == "1" {
		cfg := new(conflictless.Config)

		nonExistentDir := os.TempDir() + "/foo/bar"
		cfg.Flags.Directory = &nonExistentDir

		conflictless.Check(cfg)

		return
	}

	//nolint:gosec // this is a test package so G204 doesn't really matter here.
	cmd := exec.Command(os.Args[0], "-test.run=^TestCheckWhenDirectoryDoesNotExist$")

	cmd.Env = append(os.Environ(), "TEST_CHECK_WHEN_ERROR=1")
	err := cmd.Run()

	exitErr := new(*exec.ExitError)
	errors.As(err, exitErr)

	expectedCode := 2
	exitCode := (*exitErr).ExitCode()

	assert.Equal(t, expectedCode, exitCode, "process exited with %d, want exit status %d", expectedCode, exitCode)
}
