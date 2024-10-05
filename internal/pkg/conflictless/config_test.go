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

func TestSetChangeFileFormatFromFlags(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		description string
		format      string
		expected    string
	}{
		{"yml", "yml", "yml"},
		{"yaml", "yaml", "yaml"},
		{"json", "json", "json"},
		{"upper_case_json", "JSON", "json"},
		{"mixed_case_yaml", "yAmL", "yaml"},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			cfg := new(conflictless.Config)
			cfg.Flags.ChangeFileFormat = &testCase.format

			err := cfg.SetChangeFileFormatFromFlags()
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, cfg.ChangeFileFormat)
		})
	}
}

func TestChangeFileFormatFromFlagsWithNil(t *testing.T) {
	t.Parallel()

	cfg := new(conflictless.Config)
	cfg.ChangeFileFormat = "yml"
	cfg.Flags.ChangeFileFormat = nil

	err := cfg.SetChangeFileFormatFromFlags()
	assert.NoError(t, err)
	assert.Equal(t, "yml", cfg.ChangeFileFormat)
}

func TestChangeFileFormatFromFlagsWithInvalid(t *testing.T) {
	t.Parallel()

	invalidFormat := "foo"

	cfg := new(conflictless.Config)
	cfg.Flags.ChangeFileFormat = &invalidFormat

	err := cfg.SetChangeFileFormatFromFlags()
	assert.Error(t, err)
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
