package conflictless_test

import (
	"errors"
	"net/url"
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

func TestCliWithInvalidCommand(t *testing.T) {
	t.Parallel()

	if os.Getenv("TEST_CLI_WITH_INVALID_COMMAND") == "1" {
		conflictless.CLI()

		return
	}

	stdoutFile := createTempFile(t, os.TempDir(), "test-cli-with-invalid-command-stdout")
	defer os.Remove(stdoutFile.Name())

	stderrFile := createTempFile(t, os.TempDir(), "test-cli-with-invalid-command-stderr")
	defer os.Remove(stderrFile.Name())

	//nolint:gosec // this is a test package so G204 doesn't really matter here.
	cmd := exec.Command(os.Args[0], "-test.run=^TestCliWithInvalidCommand$", "unknown")

	cmd.Env = append(os.Environ(), "TEST_CLI_WITH_INVALID_COMMAND=1")
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile
	err := cmd.Run()

	assert.IsType(t, new(exec.ExitError), err)

	exitErr := new(*exec.ExitError)
	errors.As(err, exitErr)

	stdoutData, err := os.ReadFile(stdoutFile.Name())
	assert.NoError(t, err)
	assert.Empty(t, stdoutData)

	stderrData, err := os.ReadFile(stderrFile.Name())
	assert.NoError(t, err)
	assert.Contains(t, string(stderrData), "Error: invalid command: 'unknown'")

	expectedCode := 2
	exitCode := (*exitErr).ExitCode()

	assert.Equal(t, expectedCode, exitCode, "process exited with %d, want exit status %d", expectedCode, exitCode)
}

func TestCLIHelp(t *testing.T) {
	t.Parallel()

	if os.Getenv("TEST_CLI_HELP") != "" {
		conflictless.CLI()

		return
	}

	for _, testCase := range []struct {
		description string
		args        []string
		isError     bool
	}{
		{"help", []string{"help"}, false},
		{"help check", []string{"help", "check"}, false},
		{"help create", []string{"help", "create"}, false},
		{"help generate", []string{"help", "generate"}, false},
		{"help unknown", []string{"help", "unknown"}, true},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			stdoutFile := createTempFile(t, os.TempDir(), "test-cli-"+url.QueryEscape(testCase.description)+"-stdout")
			defer os.Remove(stdoutFile.Name())

			stderrFile := createTempFile(t, os.TempDir(), "test-cli-"+url.QueryEscape(testCase.description)+"-stderr")
			defer os.Remove(stderrFile.Name())

			//nolint:gosec // this is a test package so G204 doesn't really matter here.
			cmd := exec.Command(os.Args[0], append([]string{"-test.run=^TestCLIHelp$"}, testCase.args...)...)
			cmd.Env = append(os.Environ(), "TEST_CLI_HELP=1")
			cmd.Stdout = stdoutFile
			cmd.Stderr = stderrFile

			err := cmd.Run()

			if testCase.isError {
				assert.Error(t, err)
				assert.IsType(t, new(exec.ExitError), err)
			} else {
				assert.NoError(t, err)
			}

			stderrData, err := os.ReadFile(stderrFile.Name())
			assert.NoError(t, err)

			stdoutData, err := os.ReadFile(stdoutFile.Name())
			assert.NoError(t, err)

			if testCase.isError {
				assert.Empty(t, string(stdoutData))
				assert.Contains(t, string(stderrData), "Error:")
			} else {
				assert.Contains(t, string(stdoutData), "Usage: conflictless")
				assert.Empty(t, string(stderrData))
			}
		})
	}
}

func TestCLIGenerateWithInvalidFlags(t *testing.T) {
	t.Parallel()

	if os.Getenv("TEST_CLI_GENERATE_INVALID_FLAGS") != "" {
		conflictless.CLI()

		return
	}

	stdoutFile := createTempFile(t, os.TempDir(), "test-cli-generate-with-invalid-flags-stdout")
	defer os.Remove(stdoutFile.Name())

	stderrFile := createTempFile(t, os.TempDir(), "test-cli-generate-with-invalid-flags-stderr")
	defer os.Remove(stderrFile.Name())

	//nolint:gosec // this is a test package so G204 doesn't really matter here.
	cmd := exec.Command(
		os.Args[0],
		"-test.run=^TestCLIGenerateWithInvalidFlags$",
		"generate",
		"--dir",
		"rhymenocerous",
		"--changelog",
		"HIPPOPOTAMUS.md",
		"--bump",
		"steve",
	)

	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile

	cmd.Env = append(os.Environ(), "TEST_CLI_GENERATE_INVALID_FLAGS=1")
	err := cmd.Run()

	assert.Error(t, err)

	assert.IsType(t, new(exec.ExitError), err)

	exitErr := new(*exec.ExitError)
	errors.As(err, exitErr)

	expectedCode := 2
	exitCode := (*exitErr).ExitCode()

	assert.Equal(t, expectedCode, exitCode, "process exited with %d, want exit status %d", expectedCode, exitCode)

	stdoutData, err := os.ReadFile(stdoutFile.Name())
	assert.NoError(t, err)

	stderrData, err := os.ReadFile(stderrFile.Name())
	assert.NoError(t, err)

	assert.Empty(t, string(stdoutData))
	assert.NotEmpty(t, string(stderrData))
}

func TestCLICreateWithInvalidFlags(t *testing.T) {
	t.Parallel()

	if os.Getenv("TEST_CLI_CREATE_INVALID_FLAGS") != "" {
		conflictless.CLI()

		return
	}

	stdoutFile := createTempFile(t, os.TempDir(), "test-cli-create-with-invalid-flags-stdout")
	defer os.Remove(stdoutFile.Name())

	stderrFile := createTempFile(t, os.TempDir(), "test-cli-create-with-invalid-flags-stderr")
	defer os.Remove(stderrFile.Name())

	//nolint:gosec // this is a test package so G204 doesn't really matter here.
	cmd := exec.Command(
		os.Args[0],
		"-test.run=^TestCLICreateWithInvalidFlags$",
		"create",
		"--dir",
		"jaybird",
		"--format",
		"xml",
		"--types",
		"added,changed",
		"--name",
		"harmless-catfish",
	)

	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile

	cmd.Env = append(os.Environ(), "TEST_CLI_CREATE_INVALID_FLAGS=1")
	err := cmd.Run()

	assert.Error(t, err)
	assert.IsType(t, new(exec.ExitError), err)

	exitErr := new(*exec.ExitError)
	errors.As(err, exitErr)

	expectedCode := 2
	exitCode := (*exitErr).ExitCode()

	assert.Equal(t, expectedCode, exitCode, "process exited with %d, want exit status %d", expectedCode, exitCode)

	stdoutData, err := os.ReadFile(stdoutFile.Name())
	assert.NoError(t, err)

	stderrData, err := os.ReadFile(stderrFile.Name())
	assert.NoError(t, err)

	assert.Empty(t, string(stdoutData))
	assert.NotEmpty(t, string(stderrData))
}
