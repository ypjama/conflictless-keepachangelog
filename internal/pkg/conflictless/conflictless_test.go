package conflictless_test

import (
	"os"
	"testing"
)

const (
	mkdirFileMode = 0o755
	writeFileMode = 0o644
)

func writeDataToFile(t *testing.T, data []byte, file *os.File) {
	t.Helper()

	err := os.WriteFile(file.Name(), data, writeFileMode)
	if err != nil {
		t.Fatal(err)
	}
}

func createFile(t *testing.T, dir, name string) *os.File {
	t.Helper()

	file, err := os.Create(dir + "/" + name)
	if err != nil {
		t.Fatal(err)
	}

	return file
}

func createTempFile(t *testing.T, dir, pattern string) *os.File {
	t.Helper()

	file, err := os.CreateTemp(dir, pattern)
	if err != nil {
		t.Fatal(err)
	}

	return file
}
