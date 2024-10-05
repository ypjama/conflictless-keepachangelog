package conflictless

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/text/unicode/norm"
)

func readFile(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		return ""
	}

	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return ""
	}

	fileBytes := make([]byte, stats.Size())
	bufr := bufio.NewReader(file)

	_, err = bufr.Read(fileBytes)
	if err != nil {
		return ""
	}

	return string(fileBytes)
}

// ParseRepositoryURL detects the repository URL from the git config file.
func ParseRepositoryURL(cfg *Config) string {
	gitConfig := readFile(cfg.RepositoryConfigFile)

	re := regexp.MustCompile(`\[remote "origin"\][^[]+url = (.+)`)
	matches := re.FindStringSubmatch(gitConfig)

	if len(matches) > 1 {
		return HTTPSURLFromGitRemoteOrigin(matches[1])
	}

	return ""
}

// ParseCurrentGitBranchAsFilename parses the current git branch name from the .git/HEAD file.
func ParseCurrentGitBranchAsFilename(cfg *Config) (string, error) {
	headFile := readFile(cfg.RepositoryHeadFile)
	extension := cfg.CreateExtension

	re := regexp.MustCompile(`ref: refs/heads/(.+)`)
	matches := re.FindStringSubmatch(headFile)

	if len(matches) > 1 {
		return fmt.Sprintf("%s.%s", Basename(matches[1]), Basename(extension)), nil
	}

	return "", ErrFailedToParseBranch
}

// Basename takes a string and converts it to valid basename.
func Basename(branch string) string {
	branch = norm.NFD.String(strings.TrimSpace(branch))

	for _, reg := range []struct {
		expression  string
		replacement string
	}{
		{`[\.\/_]+`, "-"},
		{`\p{Mn}+`, ""},
		{`[-]{2}`, "-"},
	} {
		re := regexp.MustCompile(reg.expression)
		branch = re.ReplaceAllString(branch, reg.replacement)
	}

	return branch
}

// HTTPSURLFromGitRemoteOrigin converts a git remote origin URL to an HTTPS URL.
func HTTPSURLFromGitRemoteOrigin(origin string) string {
	if origin[:4] == "git@" {
		origin = origin[4:]

		re := regexp.MustCompile(`:(\w+)/`)
		origin = re.ReplaceAllString(origin, "/$1/")
	}

	for _, prefix := range []string{"https://", "http://"} {
		if origin[:len(prefix)] == prefix {
			origin = origin[len(prefix):]
		}
	}

	for _, suffix := range []string{".git", "/"} {
		if origin[len(origin)-len(suffix):] == suffix {
			origin = origin[:len(origin)-len(suffix)]
		}
	}

	return "https://" + origin
}

// ParseReleaseHeaders parses the release headers from a changelog file.
func ParseReleaseHeaders(data []byte) []string {
	headers := make([]string, 0)

	re := regexp.MustCompile(`##\s*\[(.+)]`)
	matches := re.FindAllStringSubmatch(string(data), -1)

	for _, match := range matches {
		headers = append(headers, strings.Trim(match[1], " v"))
	}

	return headers
}
