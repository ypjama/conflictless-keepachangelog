package conflictless

import "errors"

var (
	ErrChangelogFileNotFound = errors.New("changelog file not found")
	ErrChangelogWrite        = errors.New("changelog write error")
	ErrDirectoryRead         = errors.New("directory read error")
	ErrFileRead              = errors.New("file read error")
	ErrFileRemove            = errors.New("file remove error")
	ErrInvalidBumpFlag       = errors.New("invalid bump flag")
)
