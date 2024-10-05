package conflictless_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ypjama/conflictless-keepachangelog/internal/pkg/conflictless"

	"github.com/stretchr/testify/assert"
)

const (
	gitConfig = `[core]
        repositoryformatversion = 0
        filemode = true
        bare = false
        logallrefupdates = true
[remote "origin"]
        fetch = +refs/heads/*:refs/remotes/origin/*
        url = git@localhost:foo/bar-baz.git
[branch "main"]
        remote = origin
        merge = refs/heads/main
[user]
        name = foobar
        email = foobar@localhost`
)

func TestDetectRepositoryURL(t *testing.T) {
	t.Parallel()

	file, err := os.CreateTemp(os.TempDir(), "gitconfig")
	assert.NoError(t, err)

	defer os.Remove(file.Name())

	_, err = file.WriteString(gitConfig)
	assert.NoError(t, err)

	cfg := new(conflictless.Config)
	cfg.RepositoryConfigFile = file.Name()

	actual := conflictless.ParseRepositoryURL(cfg)
	assert.Equal(t, "https://localhost/foo/bar-baz", actual)
}

func TestHTTPSURLFromGitRemoteOrigin(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description string
		origin      string
		expected    string
	}

	for _, testCase := range []testCase{
		{"https", "https://github.com/golang/go.git", "https://github.com/golang/go"},
		{"http", "http://gitlab.localhost/foo/bar.git", "https://gitlab.localhost/foo/bar"},
		{"ssh", "git@github.com:golang/vscode-go.git", "https://github.com/golang/vscode-go"},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			actual := conflictless.HTTPSURLFromGitRemoteOrigin(testCase.origin)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestParseReleaseHeaders(t *testing.T) {
	t.Parallel()

	actual := conflictless.ParseReleaseHeaders([]byte(changelogContent))
	assert.Equal(t, []string{"Unreleased", "1.0.1", "1.0.0", "0.2.0", "0.1.0"}, actual)
}

func TestBasename(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description string
		input       string
		expected    string
	}

	for _, testCase := range []testCase{
		{"single ascii word", "foo", "foo"},
		{"spaces around word", " foo ", "foo"},
		{"semver", "1.0.0", "1-0-0"},
		{"semver with v", "v1.0.0", "v1-0-0"},
		{"release branch", "releases/1.0.0", "releases-1-0-0"},
		{"hotfix branch", "hotfix/did-a-thing", "hotfix-did-a-thing"},
		{"multiple-slashes", "foo/bar/-baz", "foo-bar-baz"},
		{"underscores", "foo_bar_baz", "foo-bar-baz"},
		{"underscores and dashes", "qux_quux-corge", "qux-quux-corge"},
		{"umlauts to ascii", "föö-bär-båz", "foo-bar-baz"},
		{"backwards slashes", `foo\bar\baz`, "foo-bar-baz"},
		{"kanjis are omitted", "朝日biiru", "biiru"},
		{"question mark is omitted", "foo-or-bar?", "foo-or-bar"},
		{"exclamation mark is omitted", "foo-of-course!", "foo-of-course"},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			actual := conflictless.Basename(testCase.input)
			assert.Equal(
				t,
				testCase.expected,
				actual,
				fmt.Sprintf("with input '%s' we got %s but we wanted %s", testCase.input, actual, testCase.expected),
			)
		})
	}
}
