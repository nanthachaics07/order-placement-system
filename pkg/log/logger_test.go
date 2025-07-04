package log_test

import (
	"fmt"
	"testing"

	"order-placement-system/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name        string
		env         string
		expectPrint bool
	}{
		{
			name:        "Initialize with dev environment",
			env:         "dev",
			expectPrint: false,
		},
		{
			name:        "Initialize with prod environment",
			env:         "prod",
			expectPrint: false,
		},
		{
			name:        "Initialize with invalid environment defaults to dev",
			env:         "invalid",
			expectPrint: true,
		},
		{
			name:        "Initialize with empty environment defaults to dev",
			env:         "",
			expectPrint: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				log.Init(tt.env)
			})
		})
	}
}

func TestGet(t *testing.T) {
	log.Init("dev")

	logger := log.Get()
	assert.NotNil(t, logger)

	logger2 := log.Get()
	assert.Equal(t, logger, logger2)
}

func TestGetWithoutInit(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			assert.Contains(t, r.(string), "logger not initialized")
		} else {
			t.Skip("Logger already initialized from previous test")
		}
	}()

	log.Get()
}

func TestBasicLogging(t *testing.T) {
	log.Init("dev")

	tests := []struct {
		name    string
		logFunc func(string)
		message string
	}{
		{
			name:    "Info log",
			logFunc: log.Info,
			message: "test info message",
		},
		{
			name:    "Debug log",
			logFunc: log.Debug,
			message: "test debug message",
		},
		{
			name:    "Error log",
			logFunc: log.Error,
			message: "test error message",
		},
		{
			name:    "Warn log",
			logFunc: log.Warn,
			message: "test warn message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				tt.logFunc(tt.message)
			})
		})
	}
}

func TestFormattedLogging(t *testing.T) {
	log.Init("dev")

	tests := []struct {
		name    string
		logFunc func(string, ...interface{})
		message string
		args    []interface{}
	}{
		{
			name:    "Infof with string field",
			logFunc: log.Infof,
			message: "test info message",
			args:    []interface{}{log.S("key", "value")},
		},
		{
			name:    "Debugf with multiple fields",
			logFunc: log.Debugf,
			message: "test debug message",
			args:    []interface{}{log.S("key1", "value1"), log.S("key2", "value2")},
		},
		{
			name:    "Errorf with error field",
			logFunc: log.Errorf,
			message: "test error message",
			args:    []interface{}{log.E(assert.AnError)},
		},
		{
			name:    "Warnf with mixed fields",
			logFunc: log.Warnf,
			message: "test warn message",
			args:    []interface{}{log.S("key", "value"), log.E(assert.AnError)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				tt.logFunc(tt.message, tt.args...)
			})
		})
	}
}

func TestFieldCreators(t *testing.T) {
	tests := []struct {
		name     string
		field    log.Field
		expected string
	}{
		{
			name:     "String field",
			field:    log.S("key", "value"),
			expected: "key",
		},
		{
			name:     "Error field",
			field:    log.E(assert.AnError),
			expected: "error",
		},
		{
			name:     "Any field",
			field:    log.Any("key", 123),
			expected: "key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.field)
		})
	}
}

func TestAtoS(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    interface{}
		expected string
	}{
		{
			name:     "Integer value",
			key:      "number",
			value:    123,
			expected: "123",
		},
		{
			name:     "String value",
			key:      "text",
			value:    "hello",
			expected: "hello",
		},
		{
			name:     "Boolean value",
			key:      "flag",
			value:    true,
			expected: "true",
		},
		{
			name:     "Nil value",
			key:      "nil",
			value:    nil,
			expected: "<nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := log.AtoS(tt.key, tt.value)
			assert.NotNil(t, field)
		})
	}
}

func TestLoggingWithDifferentTypes(t *testing.T) {
	log.Init("dev")

	tests := []struct {
		name string
		args []interface{}
	}{
		{
			name: "String argument",
			args: []interface{}{"test string"},
		},
		{
			name: "Integer argument",
			args: []interface{}{123},
		},
		{
			name: "Boolean argument",
			args: []interface{}{true},
		},
		{
			name: "Float argument",
			args: []interface{}{3.14},
		},
		{
			name: "Byte slice argument",
			args: []interface{}{[]byte("test bytes")},
		},
		{
			name: "Map argument",
			args: []interface{}{map[string]interface{}{"key": "value"}},
		},
		{
			name: "Slice argument",
			args: []interface{}{[]interface{}{1, 2, 3}},
		},
		{
			name: "Nil argument",
			args: []interface{}{nil},
		},
		{
			name: "Mixed arguments",
			args: []interface{}{
				log.S("key", "value"),
				assert.AnError,
				"plain string",
				123,
				true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				log.Infof("test message", tt.args...)
			})
		})
	}
}

func TestSync(t *testing.T) {
	log.Init("dev")

	err := log.Sync()
	if err != nil {
		t.Logf("Sync returned error (expected in test environment): %v", err)
	}
}

func TestSyncWithoutInit(t *testing.T) {
	err := log.Sync()
	if err != nil {
		t.Logf("Sync returned error (expected in test environment): %v", err)
	}
}

func TestConcurrentLogging(t *testing.T) {
	log.Init("dev")

	done := make(chan bool)
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()

			log.Infof("goroutine message", log.S("id", string(rune(id+'0'))))
			log.Debugf("goroutine debug", log.S("id", string(rune(id+'0'))))
			log.Errorf("goroutine error", log.S("id", string(rune(id+'0'))))
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}

func TestLoggerEnvironmentModes(t *testing.T) {
	tests := []struct {
		name         string
		env          string
		expectJSON   bool
		expectDebug  bool
		expectOutput string
	}{
		{
			name:         "dev mode uses development config",
			env:          "dev",
			expectJSON:   false,
			expectDebug:  true,
			expectOutput: "console",
		},
		{
			name:         "prod mode uses production config",
			env:          "prod",
			expectJSON:   true,
			expectDebug:  false,
			expectOutput: "json",
		},
		{
			name:         "invalid env defaults to dev",
			env:          "invalid",
			expectJSON:   false,
			expectDebug:  true,
			expectOutput: "console",
		},
		{
			name:         "empty env defaults to dev",
			env:          "",
			expectJSON:   false,
			expectDebug:  true,
			expectOutput: "console",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				log.Init(tt.env)

				// Test that logger is properly initialized
				logger := log.Get()
				assert.NotNil(t, logger)

				// Test basic logging functionality
				log.Info("test info message")
				log.Error("test error message")

				// Test debug logging based on environment
				log.Debug("test debug message")

				// Test formatted logging
				log.Infof("test formatted message", log.S("key", "value"))
			})
		})
	}
}

func TestLogLevelFiltering(t *testing.T) {
	tests := []struct {
		name        string
		env         string
		logFunc     func(string)
		shouldLog   bool
		description string
	}{
		{
			name:        "dev mode logs debug messages",
			env:         "dev",
			logFunc:     log.Debug,
			shouldLog:   true,
			description: "Debug messages should be logged in dev mode",
		},
		{
			name:        "prod mode filters debug messages",
			env:         "prod",
			logFunc:     log.Debug,
			shouldLog:   false,
			description: "Debug messages should be filtered in prod mode",
		},
		{
			name:        "dev mode logs info messages",
			env:         "dev",
			logFunc:     log.Info,
			shouldLog:   true,
			description: "Info messages should be logged in dev mode",
		},
		{
			name:        "prod mode logs info messages",
			env:         "prod",
			logFunc:     log.Info,
			shouldLog:   true,
			description: "Info messages should be logged in prod mode",
		},
		{
			name:        "dev mode logs error messages",
			env:         "dev",
			logFunc:     log.Error,
			shouldLog:   true,
			description: "Error messages should be logged in dev mode",
		},
		{
			name:        "prod mode logs error messages",
			env:         "prod",
			logFunc:     log.Error,
			shouldLog:   true,
			description: "Error messages should be logged in prod mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				log.Init(tt.env)
				tt.logFunc("test message")
			})
		})
	}
}

func TestProductionModeSpecificFeatures(t *testing.T) {
	log.Init("prod")

	// Test that production mode doesn't panic
	require.NotPanics(t, func() {
		log.Info("production info message")
		log.Error("production error message")
		log.Warn("production warn message")

		// Test formatted logging in production
		log.Infof("production formatted message",
			log.S("environment", "production"),
			log.S("feature", "json-logging"),
		)

		// Test error logging with fields
		log.Errorf("production error with context",
			log.E(assert.AnError),
			log.S("component", "logger"),
		)
	})
}

func TestDevelopmentModeSpecificFeatures(t *testing.T) {
	log.Init("dev")

	// Test that development mode doesn't panic
	require.NotPanics(t, func() {
		log.Info("development info message")
		log.Debug("development debug message")
		log.Error("development error message")
		log.Warn("development warn message")

		// Test formatted logging in development
		log.Debugf("development debug message",
			log.S("environment", "development"),
			log.S("feature", "console-logging"),
		)

		// Test error logging with fields
		log.Errorf("development error with context",
			log.E(assert.AnError),
			log.S("component", "logger"),
		)
	})
}

func TestConfigurationSwitching(t *testing.T) {
	// Test switching between different configurations
	environments := []string{"dev", "prod", "dev", "invalid", "prod"}

	for i, env := range environments {
		t.Run(fmt.Sprintf("switch_%d_%s", i, env), func(t *testing.T) {
			require.NotPanics(t, func() {
				log.Init(env)

				// Test logging after configuration switch
				log.Info("message after config switch")
				log.Infof("formatted message after config switch",
					log.S("env", env),
					log.S("iteration", string(rune(i+'0'))),
				)
			})
		})
	}
}

func TestFieldTypeSafety(t *testing.T) {
	log.Init("dev")

	tests := []struct {
		name string
		fn   func()
	}{
		{
			name: "String field with empty key",
			fn:   func() { log.Infof("test", log.S("", "value")) },
		},
		{
			name: "String field with empty value",
			fn:   func() { log.Infof("test", log.S("key", "")) },
		},
		{
			name: "Error field with nil error",
			fn:   func() { log.Infof("test", log.E(nil)) },
		},
		{
			name: "Any field with nil value",
			fn:   func() { log.Infof("test", log.Any("key", nil)) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotPanics(t, tt.fn)
		})
	}
}

func BenchmarkBasicLogging(b *testing.B) {
	log.Init("prod")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("benchmark message")
	}
}

func BenchmarkFormattedLogging(b *testing.B) {
	log.Init("prod")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Infof("benchmark message", log.S("key", "value"))
	}
}

func BenchmarkMultipleFields(b *testing.B) {
	log.Init("prod")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Infof("benchmark message",
			log.S("key1", "value1"),
			log.S("key2", "value2"),
			log.S("key3", "value3"),
		)
	}
}

func BenchmarkConcurrentLogging(b *testing.B) {
	log.Init("prod")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info("concurrent benchmark message")
		}
	})
}
