package schema

import (
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// Data for single changelog section and/or single change-file.
type Data struct {
	Added      []string `json:"added"      yaml:"added"`
	Changed    []string `json:"changed"    yaml:"changed"`
	Deprecated []string `json:"deprecated" yaml:"deprecated"`
	Removed    []string `json:"removed"    yaml:"removed"`
	Fixed      []string `json:"fixed"      yaml:"fixed"`
	Security   []string `json:"security"   yaml:"security"`
}

// JSON Schema string for validating change-files.
//
// The $id urn:uuid:<uuid> is a UUIDv5 calculated with namespace "6ba7b811-9dad-11d1-80b4-00c04fd430c8" (@url)
// and name "https://github.com/ypjama/conflictless-keepachangelog".
// `uuidgen --namespace @url --sha1 --name https://github.com/ypjama/conflictless-keepachangelog`
const jsonSchema = `{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "$id": "urn:uuid:de458349-5010-5bbc-90eb-64fd8fb5839a",
  "type": "object",
	"additionalProperties": false,
  "properties": {
    "added": {
      "type": "array",
      "minItems": 0,
      "items": {
        "type": "string"
      }
    },
    "changed": {
      "type": "array",
      "minItems": 0,
      "items": {
        "type": "string"
      }
    },
    "deprecated": {
      "type": "array",
      "minItems": 0,
      "items": {
        "type": "string"
      }
    },
    "removed": {
      "type": "array",
      "minItems": 0,
      "items": {
        "type": "string"
      }
    },
    "fixed": {
      "type": "array",
      "minItems": 0,
      "items": {
        "type": "string"
      }
    },
    "security": {
      "type": "array",
      "minItems": 0,
      "items": {
        "type": "string"
      }
    }
  },
  "required": []
}`

var (
	// ErrSchemaLoader is returned when json schema cannot be loaded properly.
	ErrSchemaLoader = errors.New("schema loader error")
	// ErrValidate is returned when data validation founds errors.
	ErrValidate = errors.New("validation error")
)

// ValidateJSON takes a JSON byte slice and validates it against the JSON Schema.
func ValidateJSON(json []byte) (bool, error) {
	schemaLoader := gojsonschema.NewStringLoader(jsonSchema)
	documentLoader := gojsonschema.NewBytesLoader(json)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrSchemaLoader, err)
	}

	if result.Valid() {
		return true, nil
	}

	return false, wrapValidationErrors(result.Errors())
}

func wrapValidationErrors(errSlice []gojsonschema.ResultError) error {
	errMsg := ""

	for _, desc := range errSlice {
		errMsg += fmt.Sprintf("- %s\n", desc)
	}

	return fmt.Errorf("%w\n\n%s", ErrValidate, errMsg)
}
