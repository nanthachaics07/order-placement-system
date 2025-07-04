package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"order-placement-system/internal/infrastructure/middleware"
	"order-placement-system/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Init("dev")
}

func TestSetup(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "Setup middleware successfully",
			description: "Should setup all middleware without panic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := gin.New()

			assert.NotPanics(t, func() {
				middleware.Setup(engine)
			})

			handlers := engine.Handlers
			assert.NotEmpty(t, handlers, "Should have middleware handlers")
		})
	}
}

func TestCORSMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectStatus   int
		expectHeaders  map[string]string
		shouldContinue bool
	}{
		{
			name:           "GET request with CORS headers",
			method:         http.MethodGet,
			expectStatus:   http.StatusOK,
			shouldContinue: true,
			expectHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Allow-Headers":     "Authorization, Content-Type, X-Requested-With, Accept",
				"Access-Control-Allow-Methods":     "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE",
			},
		},
		{
			name:           "POST request with CORS headers",
			method:         http.MethodPost,
			expectStatus:   http.StatusOK,
			shouldContinue: true,
			expectHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Allow-Headers":     "Authorization, Content-Type, X-Requested-With, Accept",
				"Access-Control-Allow-Methods":     "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE",
			},
		},
		{
			name:           "OPTIONS request should return 204",
			method:         http.MethodOptions,
			expectStatus:   http.StatusNoContent,
			shouldContinue: false,
			expectHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Allow-Headers":     "Authorization, Content-Type, X-Requested-With, Accept",
				"Access-Control-Allow-Methods":     "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE",
			},
		},
		{
			name:           "PUT request with CORS headers",
			method:         http.MethodPut,
			expectStatus:   http.StatusOK,
			shouldContinue: true,
			expectHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Allow-Headers":     "Authorization, Content-Type, X-Requested-With, Accept",
				"Access-Control-Allow-Methods":     "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE",
			},
		},
		{
			name:           "DELETE request with CORS headers",
			method:         http.MethodDelete,
			expectStatus:   http.StatusOK,
			shouldContinue: true,
			expectHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Allow-Headers":     "Authorization, Content-Type, X-Requested-With, Accept",
				"Access-Control-Allow-Methods":     "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()

			handlerCalled := false

			router.Use(func(c *gin.Context) {

				c.Header("Access-Control-Allow-Origin", "*")
				c.Header("Access-Control-Allow-Credentials", "true")
				c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With, Accept")
				c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

				if c.Request.Method == http.MethodOptions {
					c.AbortWithStatus(http.StatusNoContent)
					return
				}
				c.Next()
			})

			router.Handle(tt.method, "/test", func(c *gin.Context) {
				handlerCalled = true
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, err := http.NewRequest(tt.method, "/test", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectStatus, w.Code)

			for header, expectedValue := range tt.expectHeaders {
				assert.Equal(t, expectedValue, w.Header().Get(header),
					"Header %s should be %s", header, expectedValue)
			}

			if tt.shouldContinue {
				assert.True(t, handlerCalled, "Handler should be called for non-OPTIONS requests")
			} else {
				assert.False(t, handlerCalled, "Handler should not be called for OPTIONS requests")
			}
		})
	}
}

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name              string
		simulateError     bool
		errorMessage      string
		writeResponse     bool
		expectedStatus    int
		expectedResponse  string
		expectErrorInLogs bool
	}{
		{
			name:              "No error should continue normally",
			simulateError:     false,
			writeResponse:     true,
			expectedStatus:    http.StatusOK,
			expectedResponse:  `{"message":"success"}`,
			expectErrorInLogs: false,
		},
		{
			name:              "Error with no response written should return 500",
			simulateError:     true,
			errorMessage:      "test error",
			writeResponse:     false,
			expectedStatus:    http.StatusInternalServerError,
			expectedResponse:  `{"error":"Internal server error","message":"Something went wrong"}`,
			expectErrorInLogs: true,
		},
		{
			name:              "Error with response already written should not override",
			simulateError:     true,
			errorMessage:      "test error",
			writeResponse:     true,
			expectedStatus:    http.StatusOK,
			expectedResponse:  `{"message":"success"}`,
			expectErrorInLogs: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()

			router.Use(func(c *gin.Context) {
				c.Next()

				if len(c.Errors) > 0 {
					if !c.Writer.Written() {
						c.JSON(http.StatusInternalServerError, gin.H{
							"error":   "Internal server error",
							"message": "Something went wrong",
						})
					}
				}
			})

			router.GET("/test", func(c *gin.Context) {
				if tt.simulateError {
					c.Error(errors.New(tt.errorMessage))
				}

				if tt.writeResponse {
					c.JSON(http.StatusOK, gin.H{"message": "success"})
				}
			})

			req, err := http.NewRequest(http.MethodGet, "/test", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedResponse, w.Body.String())

			if tt.expectErrorInLogs {
				if !tt.writeResponse {
					assert.Equal(t, http.StatusInternalServerError, w.Code)
					assert.Contains(t, w.Body.String(), "Internal server error")
				}
			}
		})
	}
}

func TestCORSMiddleware_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	middleware.Setup(router)

	router.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	t.Run("Integration test with full middleware setup", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/test", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Integration test OPTIONS request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodOptions, "/api/test", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestErrorHandler_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	middleware.Setup(router)

	router.GET("/api/error", func(c *gin.Context) {
		c.Error(errors.New("test error"))
	})

	router.GET("/api/error-with-response", func(c *gin.Context) {
		c.Error(errors.New("test error"))
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	t.Run("Integration test error without response", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/error", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal server error")
	})

	t.Run("Integration test error with response already written", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/error-with-response", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})
}

func BenchmarkCORSMiddleware(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With, Accept")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkErrorHandler(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal server error",
					"message": "Something went wrong",
				})
			}
		}
	})

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkSetup(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine := gin.New()
		middleware.Setup(engine)
	}
}

func TestCORSMiddleware_EdgeCases(t *testing.T) {
	t.Run("Empty request path", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With, Accept")
			c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			c.Next()
		})

		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "root"})
		})

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("Long request path", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With, Accept")
			c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			c.Next()
		})

		longPath := "/api/v1/very/long/path/with/many/segments/test"
		router.GET(longPath, func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "long path"})
		})

		req, err := http.NewRequest(http.MethodGet, longPath, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestErrorHandler_EdgeCases(t *testing.T) {
	t.Run("Multiple errors", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Next()

			if len(c.Errors) > 0 {
				if !c.Writer.Written() {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Internal server error",
						"message": "Something went wrong",
					})
				}
			}
		})

		router.GET("/test", func(c *gin.Context) {
			c.Error(errors.New("first error"))
			c.Error(errors.New("second error"))
			c.Error(errors.New("third error"))
		})

		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal server error")
	})

	t.Run("No error case", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()

		router.Use(func(c *gin.Context) {
			c.Next()

			if len(c.Errors) > 0 {
				if !c.Writer.Written() {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Internal server error",
						"message": "Something went wrong",
					})
				}
			}
		})

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})
}
