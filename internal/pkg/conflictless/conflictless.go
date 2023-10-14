package conflictless

import "os"

func isDirectory(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}

	return info.IsDir()
}
