package conflictless_test

import (
	"errors"
	"testing"

	"github.com/ypjama/conflictless-keepachangelog/internal/pkg/conflictless"

	"github.com/stretchr/testify/assert"
)

func TestSetBumpFromFlags(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		description string
		bump        string
		expected    conflictless.Bump
	}{
		{"patch", "patch", conflictless.BumpPatch},
		{"minor", "minor", conflictless.BumpMinor},
		{"major", "major", conflictless.BumpMajor},
	} {
		// Reinitialise testCase for parallel testing.
		testCase := testCase

		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			cfg := new(conflictless.Config)
			cfg.Flags.Bump = &testCase.bump

			err := cfg.SetBumpFromFlags()
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, cfg.Bump)
		})
	}
}

func TestSetBumpFromFlagsWhenInputIsInvalid(t *testing.T) {
	t.Parallel()

	cfg := new(conflictless.Config)
	cfg.Flags.Bump = new(string)
	*cfg.Flags.Bump = "foo"

	err := cfg.SetBumpFromFlags()
	assert.Error(t, err)
	assert.True(t, errors.Is(err, conflictless.ErrInvalidBumpFlag))
}
