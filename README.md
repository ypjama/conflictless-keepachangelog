# conflictless-keepachangelog

[![go](https://github.com/ypjama/conflictless-keepachangelog/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/ypjama/conflictless-keepachangelog/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/ypjama/conflictless-keepachangelog/graph/badge.svg?token=9Q4JK8VNPO)](https://codecov.io/gh/ypjama/conflictless-keepachangelog)
[![release](https://img.shields.io/github/v/release/ypjama/conflictless-keepachangelog.svg)](https://github.com/ypjama/conflictless-keepachangelog/releases/latest)

CLI tool for generating sections to "Keep a Changelog" files without causing frustrating merge conflicts.

## Why?

[Keep a changelog format](https://keepachangelog.com/en/1.1.0/) is really great for tracking changes between versions. What is not so great is having to solve merge conflicts when project has multiple merge/pull requests and each branch has additions to the `## [Unreleased]` section. This CLI tool aims to alleviate this nuisance by introducing an alternative workflow when it comes to updating the `CHANGELOG.md` file.

## Installation

Install with go

```sh
go install github.com/ypjama/conflictless-keepachangelog/cmd/conflictless@latest
```

..or use [prebuilt binaries](https://github.com/ypjama/conflictless-keepachangelog/releases/latest).

### Verification

Use GnuPG to verify prebuild binaries.

Go to <https://github.com/ypjama/conflictless-keepachangelog/releases/latest> and
download `*_checksums.txt` and `*_checksums.txt.sig` files as well as any binary archive you wish to verify.

```sh
# Change this to match the version you downloaded.
VERSION="x.y.z"

# Import the public key from this repository.
curl --silent \
  https://raw.githubusercontent.com/ypjama/conflictless-keepachangelog/refs/heads/main/73D48E8B35873132.key \
  | gpg --import

# Verify that the signature is good on the checksums file.
# You will get a warning about the key not being certified with a trusted signature unfortunately.
gpg --verify \
  conflictless-keepachangelog_${VERSION}_checksums.txt.sig \
  conflictless-keepachangelog_${VERSION}_checksums.txt

# Compute checksums and check that they match.
sha256sum --ignore-missing --check conflictless-keepachangelog_${VERSION}_checksums.txt
```

## Usage

`conflictless help`

``` txt
Usage: conflictless <command> [flags]

The commands are:

        check           Checks that change-files are valid
        create          Creates a new change-file
        generate        Generates a version entry to changelog file
        help            Prints this help message

Use "conflictless help <topic>" for more information about that topic.
```

`conflictless help check`

``` txt
Usage: conflictless check [flags]

The flags are:

        -d, --dir
                Directory where to look for change-files (default: changes)
```

`conflictless help create`

```txt
Usage: conflictless create [flags]

The flags are:

        -d, --dir
                Directory where the change-file should be created (default: changes)
        -f, --format
                File format and extension yml/yaml/json for the change-file (default: yml)
        -t, --types
                Types of changes you want for the change-file (default: changed)

                Multiple values can be given by separating values with commas.
                Example: '--format added,changed,deprecated,removed,fixed,security'.
        -n, --name
                Name for the change-file without file extension

                If this flag is not given the name will be derived from the name of the
                current git branch you're on.
```

`conflictless help generate`

``` txt
Usage: conflictless generate [flags]

The flags are:

        -b, --bump
                Bump version patch/minor/major (default: minor)
        -c, --changelog
                Changelog file (default: CHANGELOG.md)
        -d, --dir
                Directory where to look for change-files (default: changes)
        -s, --skip-version-links
                Skip version links in changelog file (default: false)
```

## Suggested workflow

Each project should have a directory for storing unreleased changes, e.g. a directory named `changes`. In this directory developers can create _YAML_ or _JSON_ files for each merge/pull request. The filename can be freely chosen and can be derived from the branch name, e.g. `fix-broken-dependency.yml`. This way each merge/pull request would have its own `changes` file and there would not be any merge conflicts regarding the changelog.

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
