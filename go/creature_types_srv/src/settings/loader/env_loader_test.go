package settings

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvLoader_Load(t *testing.T) {
	loader := New()

	type Test struct {
		Var1 string
		Var2 int
		Var3 bool
	}

	type TestConfig struct {
		Test   Test
		Var4   string
		Test_2 struct {
			Var5 string
		}
	}

	os.Setenv("TEST_VAR1", "test_value_1")
	os.Setenv("TEST_VAR2", "123")
	os.Setenv("TEST_VAR3", "true")
	os.Setenv("VAR4", "test_value_2")
	os.Setenv("TEST_2_VAR5", "test_value_3")

	config := &TestConfig{}
	err := loader.Load(config)
	assert.NoError(t, err)
	assert.Equal(t, "test_value_1", config.Test.Var1)
	assert.Equal(t, 123, config.Test.Var2)
	assert.Equal(t, true, config.Test.Var3)
	assert.Equal(t, "test_value_2", config.Var4)
	assert.Equal(t, "test_value_3", config.Test_2.Var5)
}
