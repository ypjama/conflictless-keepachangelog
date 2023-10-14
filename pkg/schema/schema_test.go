package schema_test

import (
	"conflictless-keepachangelog/pkg/schema"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateJSONWhenInvalid(t *testing.T) {
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
		{"not a json", `foo, bar, baz`, schema.ErrSchemaLoader},
		{
			"yaml",
			`---
added:
	- foo
			`,
			schema.ErrSchemaLoader,
		},
	} {
		// Reinitialise testCase for parallel testing.
		testCase := testCase

		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			ok, err := schema.ValidateJSON([]byte(testCase.json))
			assert.False(t, ok)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, testCase.err), "expected %v, got %v", testCase.err, err)
		})
	}
}

func TestValidateJSONWhenValid(t *testing.T) {
	t.Parallel()

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

func TestValidateYAMLWhenInvalid(t *testing.T) {
	t.Parallel()

	type testCaseInvalid struct {
		description string
		yaml        string
		err         error
	}

	for _, testCase := range []testCaseInvalid{
		{
			"forbidden field",
			`---
foo: bar
			`,
			schema.ErrValidate,
		},
		{
			"added includes non-string",
			`---
added: ["foo", 123, "bar"]
			`,
			schema.ErrValidate,
		},
		{
			"changed is object instead of array",
			`---
changed: { foo: "bar" }
			`,
			schema.ErrValidate,
		},
		{"not an yaml", `foo, bar, baz`, schema.ErrYamlToJSON},
		{"unconvertable yaml", `added: { false: { true: foo } }`, schema.ErrYamlToJSON},
	} {
		// Reinitialise testCase for parallel testing.
		testCase := testCase

		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			ok, err := schema.ValidateYAML([]byte(testCase.yaml))
			assert.False(t, ok)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, testCase.err), "expected %v, got %v", testCase.err, err)
		})
	}
}

func TestValidateYamlWhenValid(t *testing.T) {
	t.Parallel()

	type testCaseValid struct {
		description string
		yaml        string
	}

	for _, testCase := range []testCaseValid{
		{"empty YAML", ``},
		{
			"Field values with square bracket syntax",
			`---
added: []
changed: []
deprecated: []
			`,
		},
		{
			"Indendation with tabs",
			`---
added:
	- foo
changed:
	- bar
			`,
		},
		{
			"Indentation uses spaces and tabs",
			`---
added:
  - foo
changed:
	- bar
			`,
		},
		{
			"Varying indentation",
			`---
added:
  - foo
changed:
    - bar
removed:
	- baz
			`,
		},
		{"Simple JSON", `{"added":["foo"]}`},
	} {
		// Reinitialise testCase for parallel testing.
		testCase := testCase

		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			ok, err := schema.ValidateYAML([]byte(testCase.yaml))
			assert.True(t, ok)
			assert.NoError(t, err)
		})
	}
}
