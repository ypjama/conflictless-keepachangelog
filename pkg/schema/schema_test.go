package schema_test

import (
	"errors"
	"testing"

	"github.com/ypjama/conflictless-keepachangelog/pkg/schema"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	t.Parallel()

	data := new(schema.Data)

	assert.True(t, data.IsEmpty())

	data.Added = []string{"foo"}
	assert.False(t, data.IsEmpty())
}

func TestParseJSONWhenInvalid(t *testing.T) {
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
		{"invisible character after json", `{"added": [""] }` + string([]rune("\u200e")), schema.ErrSchemaLoader},
		{
			"yaml",
			`---
added:
	- foo
			`,
			schema.ErrSchemaLoader,
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			data, err := schema.ParseJSON([]byte(testCase.json))
			assert.Nil(t, data)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, testCase.err), "expected %v, got %v", testCase.err, err)
		})
	}
}

func TestParseJSONWhenValid(t *testing.T) {
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
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			data, err := schema.ParseJSON([]byte(testCase.json))
			assert.NotNil(t, data)
			assert.NoError(t, err)
		})
	}
}

func TestParseYAMLWhenInvalid(t *testing.T) {
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
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			data, err := schema.ParseYAML([]byte(testCase.yaml))
			assert.Nil(t, data)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, testCase.err), "expected %v, got %v", testCase.err, err)
		})
	}
}

func TestParseYAMLWhenValid(t *testing.T) {
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
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			data, err := schema.ParseYAML([]byte(testCase.yaml))
			assert.NotNil(t, data)
			assert.NoError(t, err)
		})
	}
}
