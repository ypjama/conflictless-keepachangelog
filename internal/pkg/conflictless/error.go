package conflictless

import "errors"

var (
	ErrChangelogFileNotFound = errors.New("changelog file not found")
	ErrChangelogWrite        = errors.New("changelog write error")
	ErrCreateWrite           = errors.New("change-file write error")
	ErrDirectoryRead         = errors.New("directory read error")
	ErrFailedToParseBranch   = errors.New("failed to parse current git branch name")
	ErrFileAlreadyExists     = errors.New("file already exists")
	ErrFileRead              = errors.New("file read error")
	ErrFileRemove            = errors.New("file remove error")
	ErrInvalidBumpFlag       = errors.New("invalid bump flag")
	ErrInvalidFormatFlag     = errors.New("invalid format flag")
)
