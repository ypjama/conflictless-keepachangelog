package conflictless_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ypjama/conflictless-keepachangelog/internal/pkg/conflictless"
)

type createTestCase struct {
	description    string
	format         string
	branchName     string
	changeTypesCSV string
	contains       []string
	notContains    []string
	name           *string
}

func testCasesForCreate(t *testing.T) []createTestCase {
	t.Helper()

	name := "loving-ladybird"

	return []createTestCase{
		{
			description:    "yml_format",
			name:           nil,
			format:         "yml",
			branchName:     "foo-bar-baz",
			changeTypesCSV: "changed",
			contains:       []string{"---", "changed:\n  -"},
			notContains:    []string{"added", "deprecated", "removed", "fixed", "security"},
		},
		{
			description:    "yaml_format",
			name:           nil,
			format:         "yaml",
			branchName:     "123-create-and-fix-stuff",
			changeTypesCSV: "security,fixed,added",
			contains: []string{
				"---",
				"added:\n  -",
				"fixed:\n  -",
				"security:\n  -",
			},
			notContains: []string{"changed", "deprecated", "removed"},
		},
		{
			description:    "json_format",
			name:           nil,
			format:         "json",
			branchName:     "changing-deprecating-and-removing",
			changeTypesCSV: "changed,deprecated,removed",
			contains: []string{
				"{\n",
				"\n}",
				"\n  \"changed\": [\n    \"\"\n  ]",
				"\n  \"deprecated\": [\n    \"\"\n  ]",
				"\n  \"removed\": [\n    \"\"\n  ]",
			},
			notContains: []string{"added", "fixed", "security"},
		},
		{
			description:    "name_given",
			name:           &name,
			format:         "yml",
			branchName:     "",
			changeTypesCSV: "added",
			contains:       []string{"---", "added:\n  -"},
			notContains:    []string{"changed", "deprecated", "removed", "fixed", "security"},
		},
	}
}

func setupCreate(
	t *testing.T,
	headFileContents,
	format,
	changeTypesCSV string,
	name *string,
) (string, string, *conflictless.Config) {
	t.Helper()

	changesDir, err := os.MkdirTemp(os.TempDir(), "changes")
	assert.NoError(t, err)

	gitHeadFile := createTempFile(t, os.TempDir(), "test-generate.HEAD")
	writeDataToFile(t, []byte(headFileContents), gitHeadFile)

	cfg := new(conflictless.Config)
	cfg.RepositoryHeadFile = gitHeadFile.Name()
	cfg.Flags.ChangeTypesCsv = &changeTypesCSV
	cfg.Flags.Directory = &changesDir
	cfg.Flags.ChangeFileFormat = &format
	cfg.Flags.ChangeFileName = name

	return changesDir, gitHeadFile.Name(), cfg
}

func TestCreate(t *testing.T) {
	t.Parallel()

	for _, testCase := range testCasesForCreate(t) {
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			changesDir, gitHeadFile, cfg := setupCreate(
				t,
				`ref: refs/heads/`+testCase.branchName,
				testCase.format,
				testCase.changeTypesCSV,
				testCase.name,
			)
			defer os.RemoveAll(changesDir)
			defer os.Remove(gitHeadFile)

			conflictless.Create(cfg)

			filename := testCase.branchName + "." + testCase.format
			if testCase.name != nil {
				filename = *testCase.name + "." + testCase.format
			}

			expectedName := filepath.Join(changesDir, filename)

			file, err := os.Stat(expectedName)
			assert.NoError(t, err)
			assert.False(t, file.IsDir())

			contentBytes, err := os.ReadFile(expectedName)
			assert.NoError(t, err)

			contents := string(contentBytes)

			for _, contains := range testCase.contains {
				assert.Contains(t, contents, contains)
			}

			for _, notContains := range testCase.notContains {
				assert.NotContains(t, contents, notContains)
			}
		})
	}
}
