package load_env_test

import (
	"os"
	"strings"
	"syscall"
	"testing"

	"order-placement-system/pkg/load_env"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	err := os.Setenv(key, value)
	require.NoError(t, err)
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()
	err := os.Unsetenv(key)
	require.NoError(t, err)
}

func TestWarnIfEmpty(t *testing.T) {
	tests := []struct {
		name        string
		envName     string
		envValue    string
		setEnv      bool
		description []string
		expected    string
	}{
		{
			name:     "Environment variable exists",
			envName:  "TEST_ENV_EXISTS",
			envValue: "test_value",
			setEnv:   true,
			expected: "test_value",
		},
		{
			name:     "Environment variable does not exist",
			envName:  "TEST_ENV_NOT_EXISTS",
			setEnv:   false,
			expected: "",
		},
		{
			name:        "Environment variable does not exist with description",
			envName:     "TEST_ENV_NOT_EXISTS_DESC",
			setEnv:      false,
			description: []string{"This is a test description"},
			expected:    "",
		},
		{
			name:     "Environment variable with empty value",
			envName:  "TEST_ENV_EMPTY",
			envValue: "",
			setEnv:   true,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				setEnv(t, tt.envName, tt.envValue)
				defer unsetEnv(t, tt.envName)
			} else {
				unsetEnv(t, tt.envName)
			}

			load_env.Assert()

			var result string
			if len(tt.description) > 0 {
				result = load_env.WarnIfEmpty(tt.envName, tt.description...)
			} else {
				result = load_env.WarnIfEmpty(tt.envName)
			}

			assert.Equal(t, tt.expected, result)

			load_env.Assert()
		})
	}
}

func TestDefault(t *testing.T) {
	tests := []struct {
		name         string
		envName      string
		envValue     string
		defaultValue string
		setEnv       bool
		expected     string
	}{
		{
			name:         "Environment variable exists",
			envName:      "TEST_DEFAULT_EXISTS",
			envValue:     "existing_value",
			defaultValue: "default_value",
			setEnv:       true,
			expected:     "existing_value",
		},
		{
			name:         "Environment variable does not exist",
			envName:      "TEST_DEFAULT_NOT_EXISTS",
			defaultValue: "default_value",
			setEnv:       false,
			expected:     "default_value",
		},
		{
			name:         "Environment variable with empty value",
			envName:      "TEST_DEFAULT_EMPTY",
			envValue:     "",
			defaultValue: "default_value",
			setEnv:       true,
			expected:     "",
		},
		{
			name:         "Empty default value",
			envName:      "TEST_DEFAULT_EMPTY_DEFAULT",
			defaultValue: "",
			setEnv:       false,
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				setEnv(t, tt.envName, tt.envValue)
				defer unsetEnv(t, tt.envName)
			} else {
				unsetEnv(t, tt.envName)
			}

			result := load_env.Default(tt.envName, tt.defaultValue)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRequire(t *testing.T) {
	tests := []struct {
		name        string
		envName     string
		envValue    string
		setEnv      bool
		description []string
		expected    string
		shouldPanic bool
	}{
		{
			name:        "Environment variable exists",
			envName:     "TEST_REQUIRE_EXISTS",
			envValue:    "required_value",
			setEnv:      true,
			expected:    "required_value",
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				setEnv(t, tt.envName, tt.envValue)
				defer unsetEnv(t, tt.envName)
			} else {
				unsetEnv(t, tt.envName)
			}

			load_env.Assert()

			var result string
			if len(tt.description) > 0 {
				result = load_env.Require(tt.envName, tt.description...)
			} else {
				result = load_env.Require(tt.envName)
			}

			if tt.shouldPanic {
				assert.Panics(t, func() {
					load_env.Assert()
				})
			} else {
				assert.Equal(t, tt.expected, result)
				assert.NotPanics(t, func() {
					load_env.Assert()
				})
			}
		})
	}
}

func TestConcurrentAccess(t *testing.T) {
	t.Run("Concurrent WarnIfEmpty", func(t *testing.T) {
		setEnv(t, "TEST_CONCURRENT", "concurrent_value")
		defer unsetEnv(t, "TEST_CONCURRENT")

		load_env.Assert()

		result1 := load_env.WarnIfEmpty("TEST_CONCURRENT")
		result2 := load_env.WarnIfEmpty("TEST_CONCURRENT")

		assert.Equal(t, "concurrent_value", result1)
		assert.Equal(t, "concurrent_value", result2)

		assert.NotPanics(t, func() {
			load_env.Assert()
		})
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("Empty env name", func(t *testing.T) {
		load_env.Assert()

		result := load_env.Default("", "default")
		assert.Equal(t, "default", result)

		assert.NotPanics(t, func() {
			load_env.Assert()
		})
	})

	t.Run("Multiple descriptions", func(t *testing.T) {
		load_env.Assert()

		unsetEnv(t, "TEST_MULTI_DESC")

		load_env.WarnIfEmpty("TEST_MULTI_DESC", "First description", "Second description")

		assert.NotPanics(t, func() {
			load_env.Assert()
		})
	})

	t.Run("Nil description slice", func(t *testing.T) {
		load_env.Assert()

		unsetEnv(t, "TEST_NIL_DESC")

		result := load_env.WarnIfEmpty("TEST_NIL_DESC")
		assert.Equal(t, "", result)

		assert.NotPanics(t, func() {
			load_env.Assert()
		})
	})
}

func BenchmarkWarnIfEmpty(b *testing.B) {
	os.Setenv("BENCH_TEST", "bench_value")
	defer os.Unsetenv("BENCH_TEST")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		load_env.WarnIfEmpty("BENCH_TEST")
	}
}

func BenchmarkDefault(b *testing.B) {
	os.Setenv("BENCH_TEST", "bench_value")
	defer os.Unsetenv("BENCH_TEST")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		load_env.Default("BENCH_TEST", "default")
	}
}

func BenchmarkRequire(b *testing.B) {
	os.Setenv("BENCH_TEST", "bench_value")
	defer os.Unsetenv("BENCH_TEST")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		load_env.Require("BENCH_TEST")
	}
}

func TestSystemCallIntegration(t *testing.T) {
	testKey := "TEST_SYSCALL_INTEGRATION"
	testValue := "syscall_test_value"

	setEnv(t, testKey, testValue)
	defer unsetEnv(t, testKey)

	value, found := syscall.Getenv(testKey)
	assert.True(t, found)
	assert.Equal(t, testValue, value)

	result := load_env.Default(testKey, "default")
	assert.Equal(t, testValue, result)
}

func TestRealEnvironmentIntegration(t *testing.T) {
	pathValue := load_env.Default("PATH", "/usr/bin")
	assert.NotEmpty(t, pathValue)

	randomValue := load_env.Default("DEFINITELY_NOT_EXISTS_"+strings.Repeat("X", 10), "default")
	assert.Equal(t, "default", randomValue)
}
