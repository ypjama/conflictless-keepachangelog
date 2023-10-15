package conflictless

import "errors"

var (
	errChangelogFileNotFound = errors.New("changelog file not found")
	errChangelogWrite        = errors.New("changelog write error")
	errDirectoryRead         = errors.New("directory read error")
	errFileRead              = errors.New("file read error")
	errFileRemove            = errors.New("file remove error")
	errInvalidBumpFlag       = errors.New("invalid bump flag")
)
