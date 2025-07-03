package presenter_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"order-placement-system/internal/adapter/presenter"
	pkgErrors "order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Init("dev")
}

func TestNewOrderPresenter(t *testing.T) {
	t.Run("should create new order presenter instance", func(t *testing.T) {
		orderPresenter := presenter.NewOrderPresenter()

		assert.NotNil(t, orderPresenter)
		assert.Implements(t, (*presenter.OrderPresenter)(nil), orderPresenter)
	})
}

func TestOrderPresenter_SuccessResponse(t *testing.T) {
	tests := []struct {
		name         string
		data         interface{}
		expectedBody map[string]interface{}
	}{
		{
			name: "success response with string data",
			data: "test message",
			expectedBody: map[string]interface{}{
				"status": "success",
				"data":   "test message",
			},
		},
		{
			name: "success response with map data",
			data: map[string]interface{}{
				"id":   1,
				"name": "test product",
			},
			expectedBody: map[string]interface{}{
				"status": "success",
				"data": map[string]interface{}{
					"id":   float64(1), // JSON unmarshals numbers as float64
					"name": "test product",
				},
			},
		},
		{
			name: "success response with slice data",
			data: []string{"item1", "item2"},
			expectedBody: map[string]interface{}{
				"status": "success",
				"data":   []interface{}{"item1", "item2"},
			},
		},
		{
			name: "success response with nil data",
			data: nil,
			expectedBody: map[string]interface{}{
				"status": "success",
				"data":   nil,
			},
		},
		{
			name: "success response with empty string",
			data: "",
			expectedBody: map[string]interface{}{
				"status": "success",
				"data":   "",
			},
		},
		{
			name: "success response with number data",
			data: 42,
			expectedBody: map[string]interface{}{
				"status": "success",
				"data":   float64(42),
			},
		},
		{
			name: "success response with boolean data",
			data: true,
			expectedBody: map[string]interface{}{
				"status": "success",
				"data":   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			orderPresenter := presenter.NewOrderPresenter()

			// Act
			orderPresenter.SuccessResponse(c, tt.data)

			// Assert
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}

func TestOrderPresenter_ErrorResponse(t *testing.T) {
	tests := []struct {
		name               string
		err                error
		expectedStatusCode int
		expectedError      string
	}{
		{
			name:               "not found error",
			err:                pkgErrors.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedError:      "entity not found",
		},
		{
			name:               "invalid input error",
			err:                pkgErrors.ErrInvalidInput,
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "invalid input",
		},
		{
			name:               "already exists error",
			err:                pkgErrors.ErrAlreadyExists,
			expectedStatusCode: http.StatusConflict,
			expectedError:      "entity already exists",
		},
		{
			name:               "unauthorized error",
			err:                pkgErrors.ErrUnauthorized,
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      "unauthorized access",
		},
		{
			name:               "forbidden error",
			err:                pkgErrors.ErrForbidden,
			expectedStatusCode: http.StatusForbidden,
			expectedError:      "forbidden",
		},
		{
			name:               "conflict error",
			err:                pkgErrors.ErrConflict,
			expectedStatusCode: http.StatusConflict,
			expectedError:      "conflict",
		},
		{
			name:               "unprocessable entity error",
			err:                pkgErrors.ErrUnprocessableEntity,
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedError:      "unprocessable entity",
		},
		{
			name:               "too many requests error",
			err:                pkgErrors.ErrTooManyRequests,
			expectedStatusCode: http.StatusTooManyRequests,
			expectedError:      "too many requests",
		},
		{
			name:               "internal server error",
			err:                pkgErrors.ErrInternalServer,
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "internal server error",
		},
		{
			name:               "bad request error",
			err:                pkgErrors.ErrBadRequest,
			expectedStatusCode: http.StatusInternalServerError, // Default case
			expectedError:      "bad request",
		},
		{
			name:               "custom error - should use default case",
			err:                errors.New("custom error message"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      "custom error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			orderPresenter := presenter.NewOrderPresenter()

			// Act
			orderPresenter.ErrorResponse(c, tt.err)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			expectedBody := map[string]interface{}{
				"error": tt.expectedError,
			}
			assert.Equal(t, expectedBody, responseBody)
		})
	}
}

func TestOrderPresenter_Integration(t *testing.T) {
	t.Run("should handle success and error responses in sequence", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		orderPresenter := presenter.NewOrderPresenter()

		// Test success response
		w1 := httptest.NewRecorder()
		c1, _ := gin.CreateTestContext(w1)

		testData := map[string]interface{}{
			"message": "operation completed successfully",
			"count":   5,
		}

		orderPresenter.SuccessResponse(c1, testData)

		assert.Equal(t, http.StatusOK, w1.Code)

		var successResponse map[string]interface{}
		err := json.Unmarshal(w1.Body.Bytes(), &successResponse)
		require.NoError(t, err)

		assert.Equal(t, "success", successResponse["status"])
		assert.NotNil(t, successResponse["data"])

		// Test error response
		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)

		orderPresenter.ErrorResponse(c2, pkgErrors.ErrInvalidInput)

		assert.Equal(t, http.StatusBadRequest, w2.Code)

		var errorResponse map[string]interface{}
		err = json.Unmarshal(w2.Body.Bytes(), &errorResponse)
		require.NoError(t, err)

		assert.Equal(t, "invalid input", errorResponse["error"])
	})
}

func TestOrderPresenter_EdgeCases(t *testing.T) {
	t.Run("should handle complex nested data structures", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		orderPresenter := presenter.NewOrderPresenter()

		complexData := map[string]interface{}{
			"orders": []map[string]interface{}{
				{
					"no":        1,
					"productId": "FG0A-CLEAR-IPHONE16PROMAX",
					"qty":       2,
					"nested": map[string]interface{}{
						"level1": map[string]interface{}{
							"level2": "deep value",
						},
					},
				},
			},
			"metadata": map[string]interface{}{
				"total":     100.50,
				"currency":  "THB",
				"processed": true,
			},
		}

		orderPresenter.SuccessResponse(c, complexData)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response["status"])
		assert.NotNil(t, response["data"])

		// Verify nested structure is preserved
		data, ok := response["data"].(map[string]interface{})
		require.True(t, ok)

		orders, ok := data["orders"].([]interface{})
		require.True(t, ok)
		assert.Len(t, orders, 1)

		metadata, ok := data["metadata"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, float64(100.50), metadata["total"])
		assert.Equal(t, "THB", metadata["currency"])
		assert.Equal(t, true, metadata["processed"])
	})

	t.Run("should handle very large data structures", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		orderPresenter := presenter.NewOrderPresenter()

		// Create large data structure
		largeData := make([]map[string]interface{}, 1000)
		for i := 0; i < 1000; i++ {
			largeData[i] = map[string]interface{}{
				"id":        i,
				"productId": "FG0A-CLEAR-IPHONE16PROMAX",
				"qty":       i + 1,
			}
		}

		orderPresenter.SuccessResponse(c, largeData)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "success", response["status"])

		data, ok := response["data"].([]interface{})
		require.True(t, ok)
		assert.Len(t, data, 1000)
	})
}

// Benchmark tests
func BenchmarkOrderPresenter_SuccessResponse(b *testing.B) {
	gin.SetMode(gin.TestMode)
	orderPresenter := presenter.NewOrderPresenter()

	testData := map[string]interface{}{
		"message": "test message",
		"count":   42,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		orderPresenter.SuccessResponse(c, testData)
	}
}

func BenchmarkOrderPresenter_ErrorResponse(b *testing.B) {
	gin.SetMode(gin.TestMode)
	orderPresenter := presenter.NewOrderPresenter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		orderPresenter.ErrorResponse(c, pkgErrors.ErrInvalidInput)
	}
}

// Test helper functions
func createTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func assertJSONResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedBody map[string]interface{}) {
	t.Helper()

	assert.Equal(t, expectedStatus, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	assert.Equal(t, expectedBody, responseBody)
}
