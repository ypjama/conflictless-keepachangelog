package conflictless_test

import (
	"bytes"
	"io"
	"os"
	"testing"
)

type readWriteCapture struct {
	original   *os.File
	read       *os.File
	write      *os.File
	outChannel chan string
}

const (
	mkdirFileMode = 0o755
	writeFileMode = 0o644
)

//nolint:gochecknoglobals
var stdoutCapture *readWriteCapture

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

func startStdoutCapture(t *testing.T) {
	t.Helper()

	stdoutCapture = new(readWriteCapture)
	stdoutCapture.original = os.Stdout

	read, write, _ := os.Pipe()
	stdoutCapture.read = read
	stdoutCapture.write = write

	os.Stdout = stdoutCapture.write

	stdoutCapture.outChannel = make(chan string)

	go func() {
		var buf bytes.Buffer

		_, err := io.Copy(&buf, read)
		if err != nil {
			t.Error(err.Error())
		}

		stdoutCapture.outChannel <- buf.String()
	}()
}

func stopStdoutCapture(t *testing.T) string {
	t.Helper()

	if stdoutCapture == nil {
		t.Errorf("could not stop stdout capture because stdout capture wasn't initialized")
	}

	stdoutCapture.write.Close()
	os.Stdout = stdoutCapture.original

	output := <-stdoutCapture.outChannel

	stdoutCapture = nil

	return output
}
