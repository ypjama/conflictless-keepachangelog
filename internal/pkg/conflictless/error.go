package conflictless

import "errors"

var (
	errDirectoryRead = errors.New("directory read error")
	errFileRead      = errors.New("file read error")
)
