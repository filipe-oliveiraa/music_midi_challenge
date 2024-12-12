package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinParents(t *testing.T) {
	{
		expected := "ENV_TEST_ONE"
		actual := joinParents([]string{"Test"}, "Env", "ONE", "_", UpperCase)
		assert.Equal(t, expected, actual)
	}

	{
		expected := "flag-test-one"
		actual := joinParents([]string{"Test"}, "flag", "One", "-", LowerCase)
		assert.Equal(t, expected, actual)
	}

}
