package scopes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_All(t *testing.T) {
	assert.True(t, All()(nil), "All should always return true")
}
