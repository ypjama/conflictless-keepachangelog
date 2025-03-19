package conflictless_test

import (
	"testing"

	"github.com/ypjama/conflictless-keepachangelog/internal/pkg/conflictless"

	"github.com/stretchr/testify/assert"
)

func TestSectionLink(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description string
		baseURL     string
		sectionName string
		useVPrefix  bool
		expected    string
	}

	for _, testCase := range []testCase{
		{"empty base URL", "", "1.0.0", false, ""},
		{
			"GitHub",
			"https://github.com/olivierlacan/keep-a-changelog",
			"1.0.0",
			true,
			"[1.0.0]: https://github.com/olivierlacan/keep-a-changelog/releases/tag/v1.0.0",
		},
		{
			"GitLab",
			"https://gitlab.com/gitlab-org/gitlab",
			"16.4.0-ee",
			true,
			"[16.4.0-ee]: https://gitlab.com/gitlab-org/gitlab/-/releases/v16.4.0-ee",
		},
		{
			"Self-hosted GitLab",
			"https://gitlab.localhost/foo/bar",
			"1.2.3",
			true,
			"[1.2.3]: https://gitlab.localhost/foo/bar/-/releases/v1.2.3",
		},
		{
			"Unknown host",
			"https://example.com/foo/bar",
			"1.2.3",
			false,
			"",
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			actual := conflictless.SectionLink(testCase.baseURL, testCase.sectionName, testCase.useVPrefix)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
