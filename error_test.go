package nsqclient

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDebugError(t *testing.T) {
	err := NewDebugError("debug error")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "debug error")

	{
		err := NewDebugError(errors.New("debug error"))
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "debug error")
	}
}
