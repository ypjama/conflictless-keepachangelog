# conflictless-keepachangelog

CLI tool for generating sections to "Keep a Changelog" files without causing frustrating merge conflicts.

## Installation

```sh
go install github.com/ypjama/conflictless-keepachangelog/cmd/conflictless@latest
```

## Usage

`conflictless help`

``` txt
Usage: conflictless <command> [flags]

The commands are:

        check           Checks that change-files are valid
        generate        Generates a version entry to changelog file
        help            Prints this help message

Use "conflictless help <topic>" for more information about that topic.
```

`conflictless help generate`

``` txt
Usage: conflictless generate [flags]

The flags are:

        -b, --bump
                Bump version patch/minor/major (default: minor)
        -d, --dir
                Directory where to look for change-files (default: changes)
        -s, --skip-version-links
                Skip version links in changelog file (default: false)
```

`conflictless help check`

```txt
Usage: conflictless check [flags]

The flags are:

        -d, --dir
                Directory where to look for change-files (default: changes)
```

## Why?

[Keep a changelog format](https://keepachangelog.com/en/1.1.0/) is really great for tracking changes between versions. What is not so great is having to solve merge conflicts when project has multiple merge/pull requests and each branch has additions to the `## [Unreleased]` section. This CLI tool aims to alleviate this nuisance by introducing an alternative workflow when it comes to updating the `CHANGELOG.md` file.

## Suggested workflow

Each project should have a directory for storing unreleased changes, e.g. a directory named `changes`. In this directory developers can create _YAML_ or _JSON_ files for each merge/pull request. The filename can be freely choosen and can be derived from the branch name, e.g. `fix-broken-dependency.yml`. This way each merge/pull request would have its own `changes` file and there would not be any merge conflicts regarding the changelog.

When it's time to do a release a maintainer of the project (or pipeline automation) can generate a new version entry to the `CHANGELOG.md` file by running `conflictless generate`.

## JSON Schema for change-files

The schema for change-files is pretty simple. You can check the [JSON Schema file here](pkg/schema/jsonschema.json).

Here's an example of a valid change-file

```yml
---
added:
  - Very important feature
  - >-
    This entry is long so I'm using this different yaml syntax.
    This will all be a single line in the CHANGELOG.md file, don't worry.
    Anyway, here's important info about the addition:
    "Lorem ipsum dolor sit amet, consectetur adipiscing elit,
    sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
changed:
  - All fonts are replaced by Comic Sans!
fixed:
  - That really annoying bug
```
