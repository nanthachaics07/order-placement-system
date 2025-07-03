package errors_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	errs "order-placement-system/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorConstants(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		expectedError string
	}{
		{
			name:          "ErrNotFound should have correct message",
			err:           errs.ErrNotFound,
			expectedError: "entity not found",
		},
		{
			name:          "ErrAlreadyExists should have correct message",
			err:           errs.ErrAlreadyExists,
			expectedError: "entity already exists",
		},
		{
			name:          "ErrInvalidInput should have correct message",
			err:           errs.ErrInvalidInput,
			expectedError: "invalid input",
		},
		{
			name:          "ErrUnauthorized should have correct message",
			err:           errs.ErrUnauthorized,
			expectedError: "unauthorized access",
		},
		{
			name:          "ErrInternalServer should have correct message",
			err:           errs.ErrInternalServer,
			expectedError: "internal server error",
		},
		{
			name:          "ErrConflict should have correct message",
			err:           errs.ErrConflict,
			expectedError: "conflict",
		},
		{
			name:          "ErrForbidden should have correct message",
			err:           errs.ErrForbidden,
			expectedError: "forbidden",
		},
		{
			name:          "ErrBadRequest should have correct message",
			err:           errs.ErrBadRequest,
			expectedError: "bad request",
		},
		{
			name:          "ErrUnprocessableEntity should have correct message",
			err:           errs.ErrUnprocessableEntity,
			expectedError: "unprocessable entity",
		},
		{
			name:          "ErrTooManyRequests should have correct message",
			err:           errs.ErrTooManyRequests,
			expectedError: "too many requests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedError, tt.err.Error())
		})
	}
}

func TestMapJsonError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		inputError         error
		expectedStatusCode int
		expectedMessage    string
	}{
		{
			name:               "ErrNotFound should map to 404",
			inputError:         errs.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedMessage:    "entity not found",
		},
		{
			name:               "ErrInvalidInput should map to 400",
			inputError:         errs.ErrInvalidInput,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "invalid input",
		},
		{
			name:               "ErrAlreadyExists should map to 409",
			inputError:         errs.ErrAlreadyExists,
			expectedStatusCode: http.StatusConflict,
			expectedMessage:    "entity already exists",
		},
		{
			name:               "ErrUnprocessableEntity should map to 422",
			inputError:         errs.ErrUnprocessableEntity,
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedMessage:    "unprocessable entity",
		},
		{
			name:               "ErrUnauthorized should map to 401",
			inputError:         errs.ErrUnauthorized,
			expectedStatusCode: http.StatusUnauthorized,
			expectedMessage:    "unauthorized access",
		},
		{
			name:               "ErrForbidden should map to 403",
			inputError:         errs.ErrForbidden,
			expectedStatusCode: http.StatusForbidden,
			expectedMessage:    "forbidden",
		},
		{
			name:               "ErrConflict should map to 409",
			inputError:         errs.ErrConflict,
			expectedStatusCode: http.StatusConflict,
			expectedMessage:    "conflict",
		},
		{
			name:               "ErrTooManyRequests should map to 429",
			inputError:         errs.ErrTooManyRequests,
			expectedStatusCode: http.StatusTooManyRequests,
			expectedMessage:    "too many requests",
		},
		{
			name:               "Unknown error should map to 500",
			inputError:         errors.New("unknown error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "unknown error",
		},
		{
			name:               "Custom error should map to 500",
			inputError:         errors.New("custom business error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "custom business error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			errs.MapJsonError(c, tt.inputError)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedMessage, response["error"])
		})
	}
}

func TestMapJsonError_ResponseFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return correct JSON format", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		errs.MapJsonError(c, errs.ErrNotFound)

		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response, 1)
		assert.Contains(t, response, "error")
		assert.IsType(t, "", response["error"])
	})
}

func TestMapJsonError_MultipleConflictCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("ErrAlreadyExists and ErrConflict should both map to 409", func(t *testing.T) {
		w1 := httptest.NewRecorder()
		c1, _ := gin.CreateTestContext(w1)
		errs.MapJsonError(c1, errs.ErrAlreadyExists)
		assert.Equal(t, http.StatusConflict, w1.Code)

		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)
		errs.MapJsonError(c2, errs.ErrConflict)
		assert.Equal(t, http.StatusConflict, w2.Code)
	})
}

func BenchmarkMapJsonError(b *testing.B) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name string
		err  error
	}{
		{"KnownError", errs.ErrNotFound},
		{"UnknownError", errors.New("unknown error")},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				errs.MapJsonError(c, tc.err)
			}
		})
	}
}

func TestErrorMessages_AreNotEmpty(t *testing.T) {
	errors := []error{
		errs.ErrNotFound,
		errs.ErrAlreadyExists,
		errs.ErrInvalidInput,
		errs.ErrUnauthorized,
		errs.ErrInternalServer,
		errs.ErrConflict,
		errs.ErrForbidden,
		errs.ErrBadRequest,
		errs.ErrUnprocessableEntity,
		errs.ErrTooManyRequests,
	}

	for _, err := range errors {
		t.Run(err.Error(), func(t *testing.T) {
			assert.NotEmpty(t, err.Error(), "Error message should not be empty")
			assert.Greater(t, len(err.Error()), 0, "Error message should have length > 0")
		})
	}
}

func TestErrorMessages_AreUnique(t *testing.T) {
	errors := []error{
		errs.ErrNotFound,
		errs.ErrAlreadyExists,
		errs.ErrInvalidInput,
		errs.ErrUnauthorized,
		errs.ErrInternalServer,
		errs.ErrConflict,
		errs.ErrForbidden,
		errs.ErrBadRequest,
		errs.ErrUnprocessableEntity,
		errs.ErrTooManyRequests,
	}

	messages := make(map[string]bool)
	for _, err := range errors {
		message := err.Error()
		assert.False(t, messages[message], "Error message '%s' should be unique", message)
		messages[message] = true
	}
}
