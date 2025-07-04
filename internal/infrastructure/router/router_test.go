package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"order-placement-system/internal/infrastructure/router"
	mockHandler "order-placement-system/internal/mock/handler"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSetupHealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		expectedFields []string
	}{
		{
			name:           "Health check endpoint should return 200",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"status", "service", "version", "timestamp"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := gin.New()
			router.SetupHealthCheck(engine)

			req, err := http.NewRequest(http.MethodGet, "/health", nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()

			engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

			body := w.Body.String()
			for _, field := range tt.expectedFields {
				assert.Contains(t, body, field)
			}
		})
	}
}

func TestOrderPlacementV1Routes(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		setupMock      func(*mockHandler.OrderHandlerInterface)
	}{
		{
			name:           "POST /api/v1/orders/process should call ProcessOrders",
			method:         http.MethodPost,
			path:           "/api/v1/orders/process",
			expectedStatus: http.StatusOK,
			setupMock: func(m *mockHandler.OrderHandlerInterface) {
				m.On("ProcessOrders", mock.AnythingOfType("*gin.Context")).Return().Run(func(args mock.Arguments) {
					c := args.Get(0).(*gin.Context)
					c.JSON(http.StatusOK, gin.H{"message": "orders processed"})
				})
			},
		},
		{
			name:           "GET /api/v1/orders/process should return 404 (method not allowed)",
			method:         http.MethodGet,
			path:           "/api/v1/orders/process",
			expectedStatus: http.StatusNotFound,
			setupMock: func(m *mockHandler.OrderHandlerInterface) {
			},
		},
		{
			name:           "POST /api/v1/orders/invalid should return 404",
			method:         http.MethodPost,
			path:           "/api/v1/orders/invalid",
			expectedStatus: http.StatusNotFound,
			setupMock: func(m *mockHandler.OrderHandlerInterface) {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := gin.New()
			mockOrderHandler := mockHandler.NewOrderHandlerInterface(t)

			tt.setupMock(mockOrderHandler)

			router.OrderPlacementV1Routes(engine, mockOrderHandler)

			req, err := http.NewRequest(tt.method, tt.path, nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()

			engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestOrderPlacementV1Routes_Integration(t *testing.T) {
	t.Run("Should setup correct route groups and endpoints", func(t *testing.T) {
		engine := gin.New()
		mockOrderHandler := mockHandler.NewOrderHandlerInterface(t)

		mockOrderHandler.On("ProcessOrders", mock.AnythingOfType("*gin.Context")).Return().Maybe()

		router.OrderPlacementV1Routes(engine, mockOrderHandler)

		routes := engine.Routes()

		var foundRoute *gin.RouteInfo
		for _, route := range routes {
			if route.Path == "/api/v1/orders/process" && route.Method == "POST" {
				foundRoute = &route
				break
			}
		}

		assert.NotNil(t, foundRoute, "Should find POST /api/v1/orders/process route")
		assert.Equal(t, "POST", foundRoute.Method)
		assert.Equal(t, "/api/v1/orders/process", foundRoute.Path)
	})
}

func TestRouterEndpointResponses(t *testing.T) {
	tests := []struct {
		name           string
		setupRoutes    func(*gin.Engine, *mockHandler.OrderHandlerInterface)
		method         string
		path           string
		expectedStatus int
		setupMock      func(*mockHandler.OrderHandlerInterface)
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Health check should return proper JSON structure",
			setupRoutes: func(engine *gin.Engine, _ *mockHandler.OrderHandlerInterface) {
				router.SetupHealthCheck(engine)
			},
			method:         http.MethodGet,
			path:           "/health",
			expectedStatus: http.StatusOK,
			setupMock:      func(m *mockHandler.OrderHandlerInterface) {},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				body := w.Body.String()
				assert.Contains(t, body, `"status":"healthy"`)
				assert.Contains(t, body, `"service"`)
				assert.Contains(t, body, `"version"`)
				assert.Contains(t, body, `"timestamp"`)
			},
		},
		{
			name: "Order process endpoint should call handler",
			setupRoutes: func(engine *gin.Engine, handler *mockHandler.OrderHandlerInterface) {
				router.OrderPlacementV1Routes(engine, handler)
			},
			method:         http.MethodPost,
			path:           "/api/v1/orders/process",
			expectedStatus: http.StatusOK,
			setupMock: func(m *mockHandler.OrderHandlerInterface) {
				m.On("ProcessOrders", mock.AnythingOfType("*gin.Context")).Return().Run(func(args mock.Arguments) {
					c := args.Get(0).(*gin.Context)
					c.JSON(http.StatusOK, gin.H{"message": "orders processed"})
				})
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				body := w.Body.String()
				assert.Contains(t, body, `"message":"orders processed"`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := gin.New()
			mockOrderHandler := mockHandler.NewOrderHandlerInterface(t)

			tt.setupMock(mockOrderHandler)

			tt.setupRoutes(engine, mockOrderHandler)

			req, err := http.NewRequest(tt.method, tt.path, nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()

			engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
		})
	}
}

func TestRouterErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "Non-existent endpoint should return 404",
			method:         http.MethodGet,
			path:           "/non-existent",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Wrong method on existing endpoint should return 404",
			method:         http.MethodGet,
			path:           "/api/v1/orders/process",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid API version should return 404",
			method:         http.MethodPost,
			path:           "/api/v2/orders/process",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := gin.New()
			mockOrderHandler := mockHandler.NewOrderHandlerInterface(t)

			router.SetupHealthCheck(engine)
			router.OrderPlacementV1Routes(engine, mockOrderHandler)

			req, err := http.NewRequest(tt.method, tt.path, nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()

			engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestHealthCheckEndpointDetails(t *testing.T) {
	t.Run("Health check should return all required fields", func(t *testing.T) {
		engine := gin.New()
		router.SetupHealthCheck(engine)

		req, err := http.NewRequest(http.MethodGet, "/health", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

		body := w.Body.String()

		requiredFields := []string{
			`"status":"healthy"`,
			`"service"`,
			`"version"`,
			`"timestamp"`,
		}

		for _, field := range requiredFields {
			assert.Contains(t, body, field, "Response should contain field: %s", field)
		}
	})
}

func TestOrderProcessEndpointBehavior(t *testing.T) {
	t.Run("Process orders endpoint should accept POST requests", func(t *testing.T) {
		engine := gin.New()
		mockOrderHandler := mockHandler.NewOrderHandlerInterface(t)

		mockOrderHandler.On("ProcessOrders", mock.AnythingOfType("*gin.Context")).Return().Run(func(args mock.Arguments) {
			c := args.Get(0).(*gin.Context)
			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": "orders processed successfully",
			})
		})

		router.OrderPlacementV1Routes(engine, mockOrderHandler)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/orders/process", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "orders processed successfully")
	})

	t.Run("Process orders endpoint should reject non-POST methods", func(t *testing.T) {
		methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch}

		for _, method := range methods {
			t.Run("Method_"+method, func(t *testing.T) {
				engine := gin.New()
				mockOrderHandler := mockHandler.NewOrderHandlerInterface(t)

				router.OrderPlacementV1Routes(engine, mockOrderHandler)

				req, err := http.NewRequest(method, "/api/v1/orders/process", nil)
				assert.NoError(t, err)

				w := httptest.NewRecorder()

				engine.ServeHTTP(w, req)

				assert.Equal(t, http.StatusNotFound, w.Code, "Method %s should return 404", method)
			})
		}
	})
}

func BenchmarkSetupHealthCheck(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine := gin.New()
		router.SetupHealthCheck(engine)
	}
}

func BenchmarkOrderPlacementV1Routes(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine := gin.New()
		mockOrderHandler := &mockHandler.OrderHandlerInterface{}
		router.OrderPlacementV1Routes(engine, mockOrderHandler)
	}
}

func BenchmarkHealthCheckEndpoint(b *testing.B) {
	engine := gin.New()
	router.SetupHealthCheck(engine)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
	}
}

func createTestEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func createMockOrderHandler(t *testing.T) *mockHandler.OrderHandlerInterface {
	return mockHandler.NewOrderHandlerInterface(t)
}

func executeRequest(engine *gin.Engine, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

func TestAllRoutes(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		setupMock      func(*mockHandler.OrderHandlerInterface)
	}{
		{
			name:           "Health Check",
			method:         "GET",
			path:           "/health",
			expectedStatus: http.StatusOK,
			setupMock:      func(m *mockHandler.OrderHandlerInterface) {},
		},
		{
			name:           "Process Orders",
			method:         "POST",
			path:           "/api/v1/orders/process",
			expectedStatus: http.StatusOK,
			setupMock: func(m *mockHandler.OrderHandlerInterface) {
				m.On("ProcessOrders", mock.AnythingOfType("*gin.Context")).Return().Run(func(args mock.Arguments) {
					c := args.Get(0).(*gin.Context)
					c.JSON(http.StatusOK, gin.H{"message": "success"})
				})
			},
		},
		{
			name:           "Invalid Method on Process",
			method:         "GET",
			path:           "/api/v1/orders/process",
			expectedStatus: http.StatusNotFound,
			setupMock:      func(m *mockHandler.OrderHandlerInterface) {},
		},
		{
			name:           "Invalid Path",
			method:         "POST",
			path:           "/api/v1/invalid",
			expectedStatus: http.StatusNotFound,
			setupMock:      func(m *mockHandler.OrderHandlerInterface) {},
		},
		{
			name:           "Root Path",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusNotFound,
			setupMock:      func(m *mockHandler.OrderHandlerInterface) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			engine := createTestEngine()
			mockOrderHandler := createMockOrderHandler(t)

			tc.setupMock(mockOrderHandler)

			router.SetupHealthCheck(engine)
			router.OrderPlacementV1Routes(engine, mockOrderHandler)

			w := executeRequest(engine, tc.method, tc.path)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestRouteRegistration(t *testing.T) {
	t.Run("Should register all expected routes", func(t *testing.T) {
		engine := createTestEngine()
		mockOrderHandler := createMockOrderHandler(t)

		router.SetupHealthCheck(engine)
		router.OrderPlacementV1Routes(engine, mockOrderHandler)

		routes := engine.Routes()

		expectedRoutes := map[string]string{
			"GET /health":                 "health check endpoint",
			"POST /api/v1/orders/process": "process orders endpoint",
		}

		for expectedRoute, description := range expectedRoutes {
			found := false
			for _, route := range routes {
				routeSignature := route.Method + " " + route.Path
				if routeSignature == expectedRoute {
					found = true
					break
				}
			}
			assert.True(t, found, "Should register %s (%s)", expectedRoute, description)
		}

		assert.GreaterOrEqual(t, len(routes), len(expectedRoutes), "Should register at least %d routes", len(expectedRoutes))
	})
}
