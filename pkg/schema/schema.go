package schema

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	_ "embed"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

const (
	yamlIndent = 2
)

// Data for single changelog section and single change-file.
type Data struct {
	Added      []string `json:"added,omitempty"      yaml:"added,omitempty"`
	Changed    []string `json:"changed,omitempty"    yaml:"changed,omitempty"`
	Deprecated []string `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Removed    []string `json:"removed,omitempty"    yaml:"removed,omitempty"`
	Fixed      []string `json:"fixed,omitempty"      yaml:"fixed,omitempty"`
	Security   []string `json:"security,omitempty"   yaml:"security,omitempty"`
}

// IsEmpty returns true if all fields are empty.
func (d *Data) IsEmpty() bool {
	return len(d.Added) == 0 &&
		len(d.Changed) == 0 &&
		len(d.Deprecated) == 0 &&
		len(d.Removed) == 0 &&
		len(d.Fixed) == 0 &&
		len(d.Security) == 0
}

// ToJSON returns contents of Data as pretty print JSON string.
func (d *Data) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return "{}", fmt.Errorf("%w", err)
	}

	return string(bytes), nil
}

// ToYAML returns contents of Data as string which can used in a YAML file.
func (d *Data) ToYAML() (string, error) {
	contents := "---\n"

	var buf bytes.Buffer

	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(yamlIndent)

	err := enc.Encode(d)
	if err != nil {
		return contents, fmt.Errorf("%w", err)
	}

	contents += buf.String()

	return contents, nil
}

// JSON Schema for validating change-files.
//
// The $id urn:uuid:<uuid> is a UUIDv5 calculated with namespace "6ba7b811-9dad-11d1-80b4-00c04fd430c8" (@url)
// and name "https://github.com/ypjama/conflictless-keepachangelog".
// `uuidgen --namespace @url --sha1 --name https://github.com/ypjama/conflictless-keepachangelog`
//
//go:embed jsonschema.json
var jsonSchema []byte

var (
	// ErrSchemaLoader is returned when json schema cannot be loaded properly.
	ErrSchemaLoader = errors.New("schema loader error")
	// ErrValidate is returned when data validation founds errors.
	ErrValidate = errors.New("validation error")
	// ErrYamlToJSON is returned when yaml cannot be converted to json.
	ErrYamlToJSON = errors.New("yaml to json conversion error")
)

// ParseJSON takes a JSON byte slice and validates it against the JSON Schema.
// It returns a Data struct if the JSON is valid.
func ParseJSON(bytes []byte) (*Data, error) {
	schemaLoader := gojsonschema.NewBytesLoader(jsonSchema)
	documentLoader := gojsonschema.NewBytesLoader(bytes)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrSchemaLoader, err)
	}

	if result.Valid() {
		data := new(Data)

		err := json.Unmarshal(bytes, data)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrSchemaLoader, err)
		}

		return data, nil
	}

	return nil, wrapValidationErrors(result.Errors())
}

// ParseYAML takes a YAML byte slice and validates it against the JSON Schema.
// It returns a Data struct if the YAML is valid.
func ParseYAML(bytes []byte) (*Data, error) {
	json, err := yamlToJSON(bytes)
	if err != nil {
		return nil, err
	}

	return ParseJSON(json)
}

func yamlToJSON(b []byte) ([]byte, error) {
	data := map[string]interface{}{}

	// Trim space to avoid yaml parser error.
	s := strings.TrimSpace(string(b))

	// Replace tabs with spaces to avoid yaml parser error.
	s = strings.ReplaceAll(s, "\t", "  ")

	err := yaml.Unmarshal([]byte(s), &data)
	if err != nil {
		return []byte{}, fmt.Errorf("%w: %w", ErrYamlToJSON, err)
	}

	json, err := json.Marshal(data)
	if err != nil {
		return []byte{}, fmt.Errorf("%w: %w", ErrYamlToJSON, err)
	}

	return json, nil
}

func wrapValidationErrors(errSlice []gojsonschema.ResultError) error {
	errMsg := ""

	for _, desc := range errSlice {
		errMsg += fmt.Sprintf("- %s\n", desc)
	}

	return fmt.Errorf("%w:\n\n%s", ErrValidate, errMsg)
}
