package conflictless_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ypjama/conflictless-keepachangelog/internal/pkg/conflictless"

	"github.com/stretchr/testify/assert"
)

func TestIsJSON(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description string
		filename    string
		expected    bool
	}

	for _, testCase := range []testCase{
		{"no extension", "foo", false},
		{"json file", "foo.json", true},
		{"json file with upper case extension", "foo.JSON", true},
		{"yml file", "foo.yml", false},
		{"yaml file", "foo.yaml", false},
		{"full path to json file", "/tmp/foo/bar/baz.json", true},
		{"full path to yml file", "/tmp/foo/bar/baz.yml", false},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			actual := conflictless.IsJSONExtension(testCase.filename)

			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestEmptyData(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description          string
		changeTypesCsv       string
		expectedMinifiedJSON string
	}

	for _, testCase := range []testCase{
		{"empty csv", "", "{}"},
		{"only added", "added", `{"added":[""]}`},
		{"added and changed", "added,changed", `{"added":[""],"changed":[""]}`},
		{"added and changed with spaces", " added , changed ", `{"added":[""],"changed":[""]}`},
		{"only fixed is valid", "addeds,change,fixed,insecurity", `{"fixed":[""]}`},
		{"mixed upper and lower case", "aDdeD,chANGed", `{"added":[""],"changed":[""]}`},
		{
			"all change types",
			"added,changed,deprecated,removed,fixed,security",
			`{"added":[""],"changed":[""],"deprecated":[""],"removed":[""],"fixed":[""],"security":[""]}`,
		},
		{"kanjis", "朝日", "{}"},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			data := conflictless.EmptyData(testCase.changeTypesCsv)

			bytes, _ := json.Marshal(data)
			minifiedJSON := string(bytes)

			assert.Equal(
				t,
				testCase.expectedMinifiedJSON,
				minifiedJSON,
				fmt.Sprintf("expected %s but got %s", testCase.expectedMinifiedJSON, minifiedJSON),
			)
		})
	}
}
