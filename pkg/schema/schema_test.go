package schema_test

import (
	"conflictless-keepachangelog/pkg/schema"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateJSON(t *testing.T) {
	t.Parallel()

	ok, err := schema.ValidateJSON([]byte(`{"foo":"bar"}`))
	assert.False(t, ok)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, schema.ErrValidate))
}
