package conflictless

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ypjama/conflictless-keepachangelog/pkg/schema"
)

func validateDir(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("%w. %w", ErrDirectoryRead, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%w. %s is not a directory", ErrDirectoryRead, dir)
	}

	return nil
}

func readChangeFiles(dir string) ([]fs.DirEntry, error) {
	err := validateDir(dir)
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("%w. %w", ErrDirectoryRead, err)
	}

	changeFiles := []fs.DirEntry{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := filepath.Ext(file.Name())
		if ext == ".yml" || ext == ".yaml" || ext == ".json" {
			changeFiles = append(changeFiles, file)
		}
	}

	return changeFiles, nil
}

func scanDir(dir string) (*schema.Data, error) {
	files, err := readChangeFiles(dir)
	if err != nil {
		return nil, err
	}

	combined := new(schema.Data)

	for _, file := range files {
		filename := filepath.Join(dir, file.Name())

		fileData, err := scanFile(filename)
		if err != nil {
			return nil, fmt.Errorf("file '%s' - %w", filename, err)
		}

		addDataFromData(fileData, combined)
	}

	return combined, nil
}

func scanFile(filename string) (*schema.Data, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("%w. %w", ErrFileRead, err)
	}

	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, fmt.Errorf("%w. %w", ErrFileRead, statsErr)
	}

	fileBytes := make([]byte, stats.Size())
	bufr := bufio.NewReader(file)

	_, err = bufr.Read(fileBytes)
	if err != nil {
		return nil, fmt.Errorf("%w. %w", ErrFileRead, err)
	}

	if filepath.Ext(filename) == ".json" {
		data, err := schema.ParseJSON(fileBytes)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		return data, nil
	}

	data, err := schema.ParseYAML(fileBytes)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}

func addDataFromData(fromData, toData *schema.Data) {
	toData.Added = append(toData.Added, fromData.Added...)
	toData.Changed = append(toData.Changed, fromData.Changed...)
	toData.Deprecated = append(toData.Deprecated, fromData.Deprecated...)
	toData.Removed = append(toData.Removed, fromData.Removed...)
	toData.Fixed = append(toData.Fixed, fromData.Fixed...)
	toData.Security = append(toData.Security, fromData.Security...)
}

func removeChangeFiles(dir string) error {
	files, err := readChangeFiles(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		filename := filepath.Join(dir, file.Name())

		err := os.Remove(filename)
		if err != nil {
			return fmt.Errorf("%w. %w", ErrFileRemove, err)
		}
	}

	return nil
}
