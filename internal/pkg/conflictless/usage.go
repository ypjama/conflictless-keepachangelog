package conflictless

import "fmt"

const (
	flagIndentation            = "\t"
	flagDescriptionIndentation = "\t\t"
	flagDescriptionDir         = flagIndentation +
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

func usage() {
	//nolint:forbidigo
	fmt.Print(`Usage: conflictless <command> [flags]

The commands are:

        check           Checks that change-files are valid
        generate        Generates a version entry to changelog file
        help            Prints this help message

Use "conflictless help <topic>" for more information about that topic.

`)
}

func usageCheck() {
	//nolint:forbidigo
	fmt.Printf(`Usage: conflictless check [flags]

The flags are:

%s
`, flagDescriptionDir)
}

func usageGenerate() {
	//nolint:forbidigo
	fmt.Printf(`Usage: conflictless generate [flags]

The flags are:

%s
%s
%s
`, flagDescriptionBump, flagDescriptionDir, flagDescriptionSkipVersionLinks)
}
