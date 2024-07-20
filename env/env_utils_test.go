package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPtrSlice(t *testing.T) {
	input := "str1,str2"

	envs := toPtrSlice(input)

	assert.Equal(t, "str1", *(*envs)[0])
	assert.Equal(t, "str2", *(*envs)[1])

	t.Logf("envs0: %s", *(*envs)[0])
	t.Logf("envs1: %s", *(*envs)[1])
}

func TestNilGetEnv(t *testing.T) {
	input := ""

	nilStr := nilGetEnv(input)

	assert.Nil(t, nilStr)
}
