package main

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	t.Parallel()

	if os.Getenv("TEST_MAIN_CLI") == "1" {
		main()

		return
	}

	//nolint:gosec // this is a test function so G204 doesn't really matter here.
	cmd := exec.Command(os.Args[0], "-test.run=^Test_main$")

	cmd.Env = append(os.Environ(), "TEST_MAIN_CLI=1")
	err := cmd.Run()

	assert.IsType(t, new(exec.ExitError), err)

	exitErr := new(*exec.ExitError)
	errors.As(err, exitErr)

	assert.Equal(t, 2, (*exitErr).ExitCode(), "process ran with err %v, want exit status 2", err)
}
