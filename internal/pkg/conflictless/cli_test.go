package conflictless_test

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ypjama/conflictless-keepachangelog/internal/pkg/conflictless"
)

func TestCLIWithoutArguments(t *testing.T) {
	t.Parallel()

	if os.Getenv("TEST_CLI_WITHOUT_ARGUMENTS") == "1" {
		conflictless.CLI()

		return
	}

	stdoutFile := createTempFile(t, os.TempDir(), "test-cli-without-args-stdout")
	defer os.Remove(stdoutFile.Name())

	stderrFile := createTempFile(t, os.TempDir(), "test-cli-without-args-stderr")
	defer os.Remove(stderrFile.Name())

	//nolint:gosec // this is a test package so G204 doesn't really matter here.
	cmd := exec.Command(os.Args[0], "-test.run=^TestCLIWithoutArguments$")

	cmd.Env = append(os.Environ(), "TEST_CLI_WITHOUT_ARGUMENTS=1")
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile
	err := cmd.Run()

	assert.IsType(t, new(exec.ExitError), err)

	exitErr := new(*exec.ExitError)
	errors.As(err, exitErr)

	expectedCode := 2
	exitCode := (*exitErr).ExitCode()

	stdoutData, err := os.ReadFile(stdoutFile.Name())
	assert.NoError(t, err)
	assert.Empty(t, stdoutData)

	stderrData, err := os.ReadFile(stderrFile.Name())
	assert.NoError(t, err)
	assert.Contains(t, string(stderrData), "Usage: conflictless <command> [flags]")

	assert.Equal(t, expectedCode, exitCode, "process exited with %d, want exit status %d", expectedCode, exitCode)
}
