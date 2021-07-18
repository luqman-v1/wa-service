package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToInt(t *testing.T) {
	result := StoI("10")
	assert.Equal(t, result, 10, "failed convert to int")
}
