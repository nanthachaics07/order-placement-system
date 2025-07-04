package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"order-placement-system/env"
	"order-placement-system/internal/infrastructure/router"
	"order-placement-system/pkg/log"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	log.Init("dev")
	env.ServiceName = "test-service"
	env.AppVersion = "v1.0.0-test"

	m.Run()
}

func setupTestEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	router.SetupHealthCheck(engine)
	return engine
}

func TestHealthCheckEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "GET /health should return healthy status",
			method:         http.MethodGet,
			path:           "/health",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":  "healthy",
				"service": "test-service",
				"version": "v1.0.0-test",
			},
		},
		{
			name:           "POST /health should return method not allowed",
			method:         http.MethodPost,
			path:           "/health",
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
		{
			name:           "PUT /health should return method not allowed",
			method:         http.MethodPut,
			path:           "/health",
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
		{
			name:           "DELETE /health should return method not allowed",
			method:         http.MethodDelete,
			path:           "/health",
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := setupTestEngine()

			req, err := http.NewRequest(tt.method, tt.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedBody["status"], response["status"])
				assert.Equal(t, tt.expectedBody["service"], response["service"])
				assert.Equal(t, tt.expectedBody["version"], response["version"])

				timestamp, exists := response["timestamp"]
				assert.True(t, exists, "timestamp should be present")
				assert.NotEmpty(t, timestamp, "timestamp should not be empty")

				timestampStr, ok := timestamp.(string)
				assert.True(t, ok, "timestamp should be a string")

				_, err = time.Parse(time.RFC3339, timestampStr)
				assert.NoError(t, err, "timestamp should be in RFC3339 format")
			}
		})
	}
}

func TestHealthCheckResponseStructure(t *testing.T) {
	engine := setupTestEngine()

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	expectedFields := []string{"status", "service", "version", "timestamp"}
	for _, field := range expectedFields {
		_, exists := response[field]
		assert.True(t, exists, "field '%s' should be present", field)
	}

	assert.IsType(t, "", response["status"])
	assert.IsType(t, "", response["service"])
	assert.IsType(t, "", response["version"])
	assert.IsType(t, "", response["timestamp"])
}

func TestHealthCheckWithDifferentEnvironments(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		appVersion  string
	}{
		{
			name:        "Development environment",
			serviceName: "order-processing-dev",
			appVersion:  "v1.0.0-dev",
		},
		{
			name:        "Production environment",
			serviceName: "order-processing-prod",
			appVersion:  "v1.0.0",
		},
		{
			name:        "Staging environment",
			serviceName: "order-processing-staging",
			appVersion:  "v1.0.0-staging",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalServiceName := env.ServiceName
			originalAppVersion := env.AppVersion

			env.ServiceName = tt.serviceName
			env.AppVersion = tt.appVersion

			defer func() {
				env.ServiceName = originalServiceName
				env.AppVersion = originalAppVersion
			}()

			engine := setupTestEngine()

			req, err := http.NewRequest(http.MethodGet, "/health", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, "healthy", response["status"])
			assert.Equal(t, tt.serviceName, response["service"])
			assert.Equal(t, tt.appVersion, response["version"])
		})
	}
}

func TestHealthCheckTimestampAccuracy(t *testing.T) {
	engine := setupTestEngine()

	beforeRequest := time.Now().UTC().Truncate(time.Second)

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	afterRequest := time.Now().UTC().Add(1 * time.Second).Truncate(time.Second)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	timestampStr, ok := response["timestamp"].(string)
	require.True(t, ok, "timestamp should be a string")

	responseTime, err := time.Parse(time.RFC3339, timestampStr)
	require.NoError(t, err, "timestamp should be parseable")

	assert.True(t, responseTime.After(beforeRequest) || responseTime.Equal(beforeRequest),
		"response timestamp should be after or equal to request start time. "+
			"Before: %v, Response: %v, After: %v", beforeRequest, responseTime, afterRequest)
	assert.True(t, responseTime.Before(afterRequest) || responseTime.Equal(afterRequest),
		"response timestamp should be before or equal to request end time. "+
			"Before: %v, Response: %v, After: %v", beforeRequest, responseTime, afterRequest)

	now := time.Now().UTC()
	timeDiff := now.Sub(responseTime)
	assert.True(t, timeDiff < 5*time.Second,
		"timestamp should be recent (within 5 seconds). Time difference: %v", timeDiff)
}

func TestHealthCheckTimestampFormat(t *testing.T) {
	engine := setupTestEngine()

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	timestampStr, ok := response["timestamp"].(string)
	require.True(t, ok, "timestamp should be a string")
	require.NotEmpty(t, timestampStr, "timestamp should not be empty")

	parsedTime, err := time.Parse(time.RFC3339, timestampStr)
	require.NoError(t, err, "timestamp should be in RFC3339 format")

	now := time.Now().UTC()

	assert.False(t, parsedTime.After(now.Add(1*time.Second)),
		"timestamp should not be in the future: %v > %v", parsedTime, now)

	assert.False(t, parsedTime.Before(now.Add(-10*time.Second)),
		"timestamp should not be too old: %v < %v", parsedTime, now.Add(-10*time.Second))
}

func TestHealthCheckConcurrency(t *testing.T) {
	engine := setupTestEngine()

	concurrentRequests := 100
	done := make(chan bool, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			defer func() { done <- true }()

			req, err := http.NewRequest(http.MethodGet, "/health", nil)
			if err != nil {
				t.Errorf("Failed to create request: %v", err)
				return
			}

			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
				return
			}

			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("Failed to unmarshal response: %v", err)
				return
			}

			if response["status"] != "healthy" {
				t.Errorf("Expected status 'healthy', got %v", response["status"])
				return
			}
		}()
	}

	for i := 0; i < concurrentRequests; i++ {
		<-done
	}
}

func TestLogRoutes(t *testing.T) {
	engine := setupTestEngine()

	engine.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})
	engine.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	assert.NotPanics(t, func() {
		router.LogRoutes(engine)
	})

	routes := engine.Routes()
	assert.GreaterOrEqual(t, len(routes), 3)

	healthRouteFound := false
	for _, route := range routes {
		if route.Path == "/health" && route.Method == "GET" {
			healthRouteFound = true
			break
		}
	}
	assert.True(t, healthRouteFound, "Health check route should be registered")
}

func TestLogRoutesWithEmptyEngine(t *testing.T) {
	engine := gin.New()

	assert.NotPanics(t, func() {
		router.LogRoutes(engine)
	})
}

func BenchmarkHealthCheck(b *testing.B) {
	engine := setupTestEngine()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				b.Errorf("Expected status 200, got %d", w.Code)
			}
		}
	})
}

func BenchmarkHealthCheckWithResponseParsing(b *testing.B) {
	engine := setupTestEngine()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				b.Errorf("Expected status 200, got %d", w.Code)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				b.Errorf("Failed to unmarshal response: %v", err)
			}
		}
	})
}
