package conflictless_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ypjama/conflictless-keepachangelog/internal/pkg/conflictless"
)

func TestInitialVersion(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		description string
		bump        conflictless.Bump
		expected    string
	}{
		{"patch", conflictless.BumpPatch, "0.0.1"},
		{"minor", conflictless.BumpMinor, "0.1.0"},
		{"major", conflictless.BumpMajor, "1.0.0"},
		{"unknown", conflictless.Bump(123), "0.1.0"},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expected, testCase.bump.InitialVersion())
		})
	}
}
