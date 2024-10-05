package conflictless_test

import (
	"errors"
	"os"
	"os/exec"
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

func TestCreateWithKanjiBranchName(t *testing.T) {
	t.Parallel()

	kanjiBranchName := "朝日"

	changesDir, gitHeadFile, cfg := setupCreate(
		t,
		`ref: refs/heads/`+kanjiBranchName,
		"yml",
		"added,changed",
		nil,
	)
	defer os.RemoveAll(changesDir)
	defer os.Remove(gitHeadFile)

	if os.Getenv("TEST_CREATE_WITH_KANJI") == "1" {
		conflictless.Create(cfg)

		return
	}

	stderrFile := createTempFile(t, os.TempDir(), "test-cli-create-with-invalid-flags-stderr")
	defer os.Remove(stderrFile.Name())

	//nolint:gosec // this is a test package so G204 doesn't really matter here.
	cmd := exec.Command(os.Args[0], "-test.run=TestCreateWithKanjiBranchName")
	cmd.Env = append(os.Environ(), "TEST_CREATE_WITH_KANJI=1")
	cmd.Stderr = stderrFile
	err := cmd.Run()

	assert.Error(t, err)
	assert.IsType(t, new(exec.ExitError), err)

	exitErr := new(*exec.ExitError)
	errors.As(err, exitErr)

	expectedCode := 2
	exitCode := (*exitErr).ExitCode()

	assert.Equal(t, expectedCode, exitCode, "process exited with %d, want exit status %d", expectedCode, exitCode)

	stderrData, err := os.ReadFile(stderrFile.Name())
	assert.NoError(t, err)
	assert.Contains(t, string(stderrData), conflictless.ErrFailedToParseBranch.Error())
}
