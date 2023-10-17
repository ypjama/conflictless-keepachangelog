package conflictless_test

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ypjama/conflictless-keepachangelog/internal/pkg/conflictless"
)

func TestPrintErrorAndExit(t *testing.T) {
	t.Parallel()

	crasherEnv := os.Getenv("TEST_PRINT_ERROR_AND_EXIT_CRASHER")

	if crasherEnv != "" {
		conflictless.PrintErrorAndExit(crasherEnv, func() {})

		return
	}

	//nolint:gosec // this is a test package so G204 doesn't really matter here.
	cmd := exec.Command(os.Args[0], "-test.run=^TestPrintErrorAndExit$")

	cmd.Env = append(os.Environ(), "TEST_PRINT_ERROR_AND_EXIT_CRASHER=foobarbaz")
	err := cmd.Run()

	assert.IsType(t, new(exec.ExitError), err)

	exitErr := new(*exec.ExitError)
	errors.As(err, exitErr)

	assert.False(t, (*exitErr).Success(), "process ran with err %v, want non zero exit status", err)
}

func TestPrintUsageAndExit(t *testing.T) {
	t.Parallel()

	crasherEnv := os.Getenv("TEST_PRINT_USAGE_AND_EXIT_CRASHER")
	cfg := new(conflictless.Config)

	if crasherEnv != "" {
		if crasherEnv != "no-cmd" {
			cfg.Flags.Command = crasherEnv
		}

		conflictless.PrintUsageAndExit(cfg)

		return
	}

	for _, crasher := range []string{
		"no-cmd",
		"check",
		"generate",
		"usage",
	} {
		// Reinitialise crasher for parallel testing.
		crasher := crasher

		t.Run(crasher, func(t *testing.T) {
			t.Parallel()

			//nolint:gosec // this is a test package so G204 doesn't really matter here.
			cmd := exec.Command(os.Args[0], "-test.run=^TestPrintUsageAndExit$")

			cmd.Env = append(os.Environ(), "TEST_PRINT_USAGE_AND_EXIT_CRASHER="+crasher)
			err := cmd.Run()

			assert.IsType(t, new(exec.ExitError), err)

			exitErr := new(*exec.ExitError)
			errors.As(err, exitErr)

			assert.False(t, (*exitErr).Success(), "process ran with err %v, want non zero exit status", err)
		})
	}
}

func TestPrintCheckSuccess(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		description string
		noContent   bool
		expected    string
	}{
		{"no content", true, "No changes found!\n"},
		{"content", false, "Change files are valid!\n"},
	} {
		// Reinitialise testCase for parallel testing.
		testCase := testCase

		t.Run("", func(t *testing.T) {
			t.Parallel()

			file := createTempFile(t, os.TempDir(), "test-stdout")
			defer os.Remove(file.Name())

			os.Stdout = file

			conflictless.PrintCheckSuccess(testCase.noContent)

			data, err := os.ReadFile(file.Name())

			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, string(data))
		})
	}
}
