package schema_test

import (
	"conflictless-keepachangelog/pkg/schema"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateJSON(t *testing.T) {
	t.Parallel()

	type testCaseInvalid struct {
		description string
		json        string
		err         error
	}

	for _, testCase := range []testCaseInvalid{
		{"forbidden field", `{"foo":"bar"}`, schema.ErrValidate},
		{"added includes non-string", `{"added":["foo", 123, "bar"]}`, schema.ErrValidate},
		{"changed is object instead of array", `{"changed":{"foo":"bar"}}`, schema.ErrValidate},
	} {
		// Reinitialise testCase for parallel testing.
		testCase := testCase

		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			ok, err := schema.ValidateJSON([]byte(testCase.json))
			assert.False(t, ok)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, testCase.err))
		})
	}

	type testCaseValid struct {
		description string
		json        string
	}

	for _, testCase := range []testCaseValid{
		{"Empty JSON", `{}`},
		{"Fields have zero lines", `{"added":[],"changed":[],"deprecated":[],"removed":[],"fixed":[],"security":[]}`},
		{
			"Each field has one line",
			`{
				"added":["foo"],
				"changed":["bar"],
				"deprecated":["baz"],
				"removed":["qux"],
				"fixed":["quux"],
				"security":["corge"]
			}`,
		},
	} {
		// Reinitialise testCase for parallel testing.
		testCase := testCase

		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			ok, err := schema.ValidateJSON([]byte(testCase.json))
			assert.True(t, ok)
			assert.NoError(t, err)
		})
	}
}
