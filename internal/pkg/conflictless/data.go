package conflictless

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ypjama/conflictless-keepachangelog/pkg/schema"
)

const (
	validExtensionJSON = ".json"
	validExtensionYAML = ".yaml"
	validExtensionYML  = ".yml"
)

func validateAndSanitizeDir(dir string) (string, error) {
	dir = strings.TrimSpace(dir)
	dir = strings.TrimSuffix(dir, string(os.PathSeparator))

	info, err := os.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("%w. %w", ErrDirectoryRead, err)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("%w. %s is not a directory", ErrDirectoryRead, dir)
	}

	return dir, nil
}

func readChangeFiles(dir string) ([]fs.DirEntry, error) {
	dir, err := validateAndSanitizeDir(dir)
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
		if ext == validExtensionJSON || ext == validExtensionYAML || ext == validExtensionYML {
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

// EmptyData returns *schema.Data with single empty string on each of the defined change types.
func EmptyData(changeTypesCsv string) *schema.Data {
	inputMap := map[string][]string{}

	for _, changeType := range strings.Split(changeTypesCsv, ",") {
		inputMap[strings.ToLower(strings.TrimSpace(changeType))] = []string{""}
	}

	data := new(schema.Data)

	bytes, err := json.Marshal(inputMap)
	if err != nil {
		return data
	}

	err = json.Unmarshal(bytes, data)
	if err != nil {
		return data
	}

	return data
}

// IsJSONExtension takes a filename and returns true if file extension is json.
func IsJSONExtension(filename string) bool {
	extLen := len(validExtensionJSON)
	lower := strings.ToLower(strings.TrimSpace(filename))

	return len(lower) > extLen && lower[len(lower)-extLen:] == validExtensionJSON
}

func createChangeFile(cfg *Config) error {
	dir, err := validateAndSanitizeDir(cfg.Directory)
	if err != nil {
		return err
	}

	name := filepath.Join(dir, cfg.ChangeFile)

	if _, err := os.Stat(name); err == nil {
		return fmt.Errorf("%w: %s", ErrFileAlreadyExists, name)
	}

	data := EmptyData(cfg.ChangeTypesCsv)

	var contents string

	if IsJSONExtension(name) {
		contents, err = data.ToJSON()
		if err != nil {
			return fmt.Errorf("%w. %w", ErrCreateWrite, err)
		}
	} else {
		contents, err = data.ToYAML()
		if err != nil {
			return fmt.Errorf("%w. %w", ErrCreateWrite, err)
		}
	}

	err = os.WriteFile(name, []byte(contents), fs.FileMode(writeFileMode))
	if err != nil {
		return fmt.Errorf("%w. %w", ErrCreateWrite, err)
	}

	return nil
}
