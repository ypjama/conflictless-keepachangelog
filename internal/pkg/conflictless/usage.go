package conflictless

import (
	"fmt"
	"os"
)

const (
	flagIndentation            = "\t"
	flagDescriptionIndentation = "\t\t"
	flagDescriptionChangelog   = flagIndentation +
		"-c, --changelog\n" +
		flagDescriptionIndentation +
		"Changelog file (default: CHANGELOG.md)"
	flagDescriptionDir = flagIndentation +
		"-d, --dir\n" +
		flagDescriptionIndentation +
		"Directory where to look for change-files (default: changes)"
	flagDescriptionBump = flagIndentation +
		"-b, --bump\n" +
		flagDescriptionIndentation +
		"Bump version patch/minor/major (default: minor)"
	flagDescriptionSkipVersionLinks = flagIndentation +
		"-s, --skip-version-links\n" +
		flagDescriptionIndentation +
		"Skip version links in changelog file (default: false)"
)

func usageText() string {
	return `Usage: conflictless <command> [flags]

The commands are:

        check           Checks that change-files are valid
        generate        Generates a version entry to changelog file
        help            Prints this help message

Use "conflictless help <topic>" for more information about that topic.

`
}

func usageTextForGenerate() string {
	return fmt.Sprintf(`Usage: conflictless generate [flags]

The flags are:

%s
%s
%s
%s
`,
		flagDescriptionBump,
		flagDescriptionChangelog,
		flagDescriptionDir,
		flagDescriptionSkipVersionLinks,
	)
}

func usageTextForCheck() string {
	return fmt.Sprintf(`Usage: conflictless check [flags]

The flags are:

%s
`,
		flagDescriptionDir,
	)
}

func usageTextForCreate() string {
	return fmt.Sprintf(`Usage: conflictless create [flags]

The flags are:

%s
`,
		flagDescriptionDir,
	)
}

func usage() {
	fmt.Fprint(os.Stdout, usageText())
}

func usageOnError() {
	fmt.Fprint(os.Stderr, usageText())
}

func usageCheck() {
	fmt.Fprint(os.Stdout, usageTextForCheck())
}

func usageCheckOnError() {
	fmt.Fprint(os.Stderr, usageTextForCheck())
}

func usageCreate() {
	fmt.Fprint(os.Stdout, usageTextForCreate())
}

func usageCreateOnError() {
	fmt.Fprint(os.Stderr, usageTextForCreate())
}

func usageGenerate() {
	fmt.Fprint(os.Stdout, usageTextForGenerate())
}

func usageGenerateOnError() {
	fmt.Fprint(os.Stderr, usageTextForGenerate())
}
