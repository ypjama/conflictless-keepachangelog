# Changelog
<!-- markdownlint-disable MD024 -->
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2024-10-06

[0.3.0]: https://github.com/ypjama/conflictless-keepachangelog/releases/tag/v0.3.0

### Added

- New command 'preview' which prints a preview of the next changelog entry.

## [0.2.0] - 2024-10-05

[0.2.0]: https://github.com/ypjama/conflictless-keepachangelog/releases/tag/v0.2.0

### Added

- New command 'create'
- New flag for "generate" command: "-c, --changelog" for overriding the default changelog file name

### Changed

- Errors are outputted to stderr instead of stdout

## [0.1.1] - 2023-10-16

[0.1.1]: https://github.com/ypjama/conflictless-keepachangelog/releases/tag/v0.1.1

### Added

- Automatic release building with goreleaser
- Golangci-lint step to go workflow
- Codecov step to go workflow

## [0.1.0] - 2023-10-15

[0.1.0]: https://github.com/ypjama/conflictless-keepachangelog/releases/tag/v0.1.0

### Added

- CLI tool for validating YAML and JSON files containing change information
- CLI tool for generating version sections to CHANGELOG.md files
