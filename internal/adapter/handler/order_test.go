package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"

	"order-placement-system/internal/adapter/handler"
	"order-placement-system/internal/adapter/handler/model"
	"order-placement-system/internal/domain/entity"
	"order-placement-system/internal/domain/value_object"
	errs "order-placement-system/pkg/errors"
	"order-placement-system/pkg/log"
)

func init() {
	log.Init("dev")
}

type MockOrderProcessor struct {
	mock.Mock
}

func (m *MockOrderProcessor) ProcessOrders(inputOrders []*entity.InputOrder) ([]*entity.CleanedOrder, error) {
	args := m.Called(inputOrders)
	return args.Get(0).([]*entity.CleanedOrder), args.Error(1)
}

type MockPresenter struct {
	mock.Mock
}

func (m *MockPresenter) SuccessResponse(c *gin.Context, data interface{}) {
	m.Called(c, data)
}

func (m *MockPresenter) ErrorResponse(c *gin.Context, err error) {
	m.Called(c, err)
}

func TestOrderHandler_ProcessOrders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Case 1: Only one product", func(t *testing.T) {
		mockProcessor := new(MockOrderProcessor)
		mockPresenter := new(MockPresenter)

		handler := handler.NewOrderHandler(mockProcessor, mockPresenter)

		inputData := []*model.InputOrder{
			{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         50.0,
				TotalPrice:        100.0,
			},
		}

		expectedResult := []*entity.CleanedOrder{
			{
				No:         1,
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(50.0),
				TotalPrice: value_object.MustNewPrice(100.0),
			},
			{
				No:         2,
				ProductId:  "WIPING-CLOTH",
				Qty:        2,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
			{
				No:         3,
				ProductId:  "CLEAR-CLEANNER",
				Qty:        2,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
		}

		mockProcessor.On("ProcessOrders", mock.MatchedBy(func(orders []*entity.InputOrder) bool {
			return len(orders) == 1 && orders[0].PlatformProductId == "FG0A-CLEAR-IPHONE16PROMAX"
		})).Return(expectedResult, nil)

		mockPresenter.On("SuccessResponse", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("[]*model.CleanedOrder")).Return()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody, _ := json.Marshal(inputData)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orders/process", bytes.NewBuffer(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessOrders(c)

		mockProcessor.AssertExpectations(t)
		mockPresenter.AssertExpectations(t)
	})

	t.Run("Case 2: One product with wrong prefix", func(t *testing.T) {
		mockProcessor := new(MockOrderProcessor)
		mockPresenter := new(MockPresenter)

		handler := handler.NewOrderHandler(mockProcessor, mockPresenter)

		inputData := []*model.InputOrder{
			{
				No:                1,
				PlatformProductId: "x2-3&FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         50.0,
				TotalPrice:        100.0,
			},
		}

		expectedResult := []*entity.CleanedOrder{
			{
				No:         1,
				ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "IPHONE16PROMAX",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(50.0),
				TotalPrice: value_object.MustNewPrice(100.0),
			},
			{
				No:         2,
				ProductId:  "WIPING-CLOTH",
				Qty:        2,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
			{
				No:         3,
				ProductId:  "CLEAR-CLEANNER",
				Qty:        2,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
		}

		mockProcessor.On("ProcessOrders", mock.MatchedBy(func(orders []*entity.InputOrder) bool {
			return len(orders) == 1 && orders[0].PlatformProductId == "x2-3&FG0A-CLEAR-IPHONE16PROMAX"
		})).Return(expectedResult, nil)

		mockPresenter.On("SuccessResponse", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("[]*model.CleanedOrder")).Return()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody, _ := json.Marshal(inputData)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orders/process", bytes.NewBuffer(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessOrders(c)

		mockProcessor.AssertExpectations(t)
		mockPresenter.AssertExpectations(t)
	})

	t.Run("Case 3: One product with wrong prefix and * symbol", func(t *testing.T) {
		mockProcessor := new(MockOrderProcessor)
		mockPresenter := new(MockPresenter)

		handler := handler.NewOrderHandler(mockProcessor, mockPresenter)

		inputData := []*model.InputOrder{
			{
				No:                1,
				PlatformProductId: "x2-3&FG0A-MATTE-IPHONE16PROMAX*3",
				Qty:               1,
				UnitPrice:         90.0,
				TotalPrice:        90.0,
			},
		}

		expectedResult := []*entity.CleanedOrder{
			{
				No:         1,
				ProductId:  "FG0A-MATTE-IPHONE16PROMAX",
				MaterialId: "FG0A-MATTE",
				ModelId:    "IPHONE16PROMAX",
				Qty:        3,
				UnitPrice:  value_object.MustNewPrice(30.0),
				TotalPrice: value_object.MustNewPrice(90.0),
			},
			{
				No:         2,
				ProductId:  "WIPING-CLOTH",
				Qty:        3,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
			{
				No:         3,
				ProductId:  "MATTE-CLEANNER",
				Qty:        3,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
		}

		mockProcessor.On("ProcessOrders", mock.MatchedBy(func(orders []*entity.InputOrder) bool {
			return len(orders) == 1 && orders[0].PlatformProductId == "x2-3&FG0A-MATTE-IPHONE16PROMAX*3"
		})).Return(expectedResult, nil)

		mockPresenter.On("SuccessResponse", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("[]*model.CleanedOrder")).Return()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody, _ := json.Marshal(inputData)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orders/process", bytes.NewBuffer(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessOrders(c)

		mockProcessor.AssertExpectations(t)
		mockPresenter.AssertExpectations(t)
	})

	t.Run("Case 4: Bundle product with / symbol", func(t *testing.T) {
		mockProcessor := new(MockOrderProcessor)
		mockPresenter := new(MockPresenter)

		handler := handler.NewOrderHandler(mockProcessor, mockPresenter)

		inputData := []*model.InputOrder{
			{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B",
				Qty:               1,
				UnitPrice:         80.0,
				TotalPrice:        80.0,
			},
		}

		expectedResult := []*entity.CleanedOrder{
			{
				No:         1,
				ProductId:  "FG0A-CLEAR-OPPOA3",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "OPPOA3",
				Qty:        1,
				UnitPrice:  value_object.MustNewPrice(40.0),
				TotalPrice: value_object.MustNewPrice(40.0),
			},
			{
				No:         2,
				ProductId:  "FG0A-CLEAR-OPPOA3-B",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "OPPOA3-B",
				Qty:        1,
				UnitPrice:  value_object.MustNewPrice(40.0),
				TotalPrice: value_object.MustNewPrice(40.0),
			},
			{
				No:         3,
				ProductId:  "WIPING-CLOTH",
				Qty:        2,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
			{
				No:         4,
				ProductId:  "CLEAR-CLEANNER",
				Qty:        2,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
		}

		mockProcessor.On("ProcessOrders", mock.MatchedBy(func(orders []*entity.InputOrder) bool {
			return len(orders) == 1 && orders[0].PlatformProductId == "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B"
		})).Return(expectedResult, nil)

		mockPresenter.On("SuccessResponse", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("[]*model.CleanedOrder")).Return()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody, _ := json.Marshal(inputData)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orders/process", bytes.NewBuffer(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessOrders(c)

		mockProcessor.AssertExpectations(t)
		mockPresenter.AssertExpectations(t)
	})

	t.Run("Case 5: Bundle product with three products", func(t *testing.T) {
		mockProcessor := new(MockOrderProcessor)
		mockPresenter := new(MockPresenter)

		handler := handler.NewOrderHandler(mockProcessor, mockPresenter)

		inputData := []*model.InputOrder{
			{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B/FG0A-MATTE-OPPOA3",
				Qty:               1,
				UnitPrice:         120.0,
				TotalPrice:        120.0,
			},
		}

		expectedResult := []*entity.CleanedOrder{
			{
				No:         1,
				ProductId:  "FG0A-CLEAR-OPPOA3",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "OPPOA3",
				Qty:        1,
				UnitPrice:  value_object.MustNewPrice(40.0),
				TotalPrice: value_object.MustNewPrice(40.0),
			},
			{
				No:         2,
				ProductId:  "FG0A-CLEAR-OPPOA3-B",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "OPPOA3-B",
				Qty:        1,
				UnitPrice:  value_object.MustNewPrice(40.0),
				TotalPrice: value_object.MustNewPrice(40.0),
			},
			{
				No:         3,
				ProductId:  "FG0A-MATTE-OPPOA3",
				MaterialId: "FG0A-MATTE",
				ModelId:    "OPPOA3",
				Qty:        1,
				UnitPrice:  value_object.MustNewPrice(40.0),
				TotalPrice: value_object.MustNewPrice(40.0),
			},
			{
				No:         4,
				ProductId:  "WIPING-CLOTH",
				Qty:        3,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
			{
				No:         5,
				ProductId:  "CLEAR-CLEANNER",
				Qty:        2,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
			{
				No:         6,
				ProductId:  "MATTE-CLEANNER",
				Qty:        1,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
		}

		mockProcessor.On("ProcessOrders", mock.MatchedBy(func(orders []*entity.InputOrder) bool {
			return len(orders) == 1 && orders[0].PlatformProductId == "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B/FG0A-MATTE-OPPOA3"
		})).Return(expectedResult, nil)

		mockPresenter.On("SuccessResponse", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("[]*model.CleanedOrder")).Return()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody, _ := json.Marshal(inputData)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orders/process", bytes.NewBuffer(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessOrders(c)

		mockProcessor.AssertExpectations(t)
		mockPresenter.AssertExpectations(t)
	})

	t.Run("Case 6: Bundle product with / and * symbols", func(t *testing.T) {
		mockProcessor := new(MockOrderProcessor)
		mockPresenter := new(MockPresenter)

		handler := handler.NewOrderHandler(mockProcessor, mockPresenter)

		inputData := []*model.InputOrder{
			{
				No:                1,
				PlatformProductId: "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3",
				Qty:               1,
				UnitPrice:         120.0,
				TotalPrice:        120.0,
			},
		}

		expectedResult := []*entity.CleanedOrder{
			{
				No:         1,
				ProductId:  "FG0A-CLEAR-OPPOA3",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "OPPOA3",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(40.0),
				TotalPrice: value_object.MustNewPrice(80.0),
			},
			{
				No:         2,
				ProductId:  "FG0A-MATTE-OPPOA3",
				MaterialId: "FG0A-MATTE",
				ModelId:    "OPPOA3",
				Qty:        1,
				UnitPrice:  value_object.MustNewPrice(40.0),
				TotalPrice: value_object.MustNewPrice(40.0),
			},
			{
				No:         3,
				ProductId:  "WIPING-CLOTH",
				Qty:        3,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
			{
				No:         4,
				ProductId:  "CLEAR-CLEANNER",
				Qty:        2,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
			{
				No:         5,
				ProductId:  "MATTE-CLEANNER",
				Qty:        1,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
		}

		mockProcessor.On("ProcessOrders", mock.MatchedBy(func(orders []*entity.InputOrder) bool {
			return len(orders) == 1 && orders[0].PlatformProductId == "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3"
		})).Return(expectedResult, nil)

		mockPresenter.On("SuccessResponse", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("[]*model.CleanedOrder")).Return()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody, _ := json.Marshal(inputData)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orders/process", bytes.NewBuffer(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessOrders(c)

		mockProcessor.AssertExpectations(t)
		mockPresenter.AssertExpectations(t)
	})

	t.Run("Case 7: Multiple products with complex combinations", func(t *testing.T) {
		mockProcessor := new(MockOrderProcessor)
		mockPresenter := new(MockPresenter)

		handler := handler.NewOrderHandler(mockProcessor, mockPresenter)

		inputData := []*model.InputOrder{
			{
				No:                1,
				PlatformProductId: "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3*2",
				Qty:               1,
				UnitPrice:         160.0,
				TotalPrice:        160.0,
			},
			{
				No:                2,
				PlatformProductId: "FG0A-PRIVACY-IPHONE16PROMAX",
				Qty:               1,
				UnitPrice:         50.0,
				TotalPrice:        50.0,
			},
		}

		expectedResult := []*entity.CleanedOrder{
			{
				No:         1,
				ProductId:  "FG0A-CLEAR-OPPOA3",
				MaterialId: "FG0A-CLEAR",
				ModelId:    "OPPOA3",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(40.0),
				TotalPrice: value_object.MustNewPrice(80.0),
			},
			{
				No:         2,
				ProductId:  "FG0A-MATTE-OPPOA3",
				MaterialId: "FG0A-MATTE",
				ModelId:    "OPPOA3",
				Qty:        2,
				UnitPrice:  value_object.MustNewPrice(40.0),
				TotalPrice: value_object.MustNewPrice(80.0),
			},
			{
				No:         3,
				ProductId:  "FG0A-PRIVACY-IPHONE16PROMAX",
				MaterialId: "FG0A-PRIVACY",
				ModelId:    "IPHONE16PROMAX",
				Qty:        1,
				UnitPrice:  value_object.MustNewPrice(50.0),
				TotalPrice: value_object.MustNewPrice(50.0),
			},
			{
				No:         4,
				ProductId:  "WIPING-CLOTH",
				Qty:        5,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
			{
				No:         5,
				ProductId:  "CLEAR-CLEANNER",
				Qty:        2,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
			{
				No:         6,
				ProductId:  "MATTE-CLEANNER",
				Qty:        2,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
			{
				No:         7,
				ProductId:  "PRIVACY-CLEANNER",
				Qty:        1,
				UnitPrice:  value_object.ZeroPrice(),
				TotalPrice: value_object.ZeroPrice(),
			},
		}

		mockProcessor.On("ProcessOrders", mock.MatchedBy(func(orders []*entity.InputOrder) bool {
			return len(orders) == 2
		})).Return(expectedResult, nil)

		mockPresenter.On("SuccessResponse", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("[]*model.CleanedOrder")).Return()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody, _ := json.Marshal(inputData)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orders/process", bytes.NewBuffer(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessOrders(c)

		mockProcessor.AssertExpectations(t)
		mockPresenter.AssertExpectations(t)
	})
}

func TestOrderHandler_ProcessOrders_ErrorCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Invalid JSON request body", func(t *testing.T) {
		mockProcessor := new(MockOrderProcessor)
		mockPresenter := new(MockPresenter)

		handler := handler.NewOrderHandler(mockProcessor, mockPresenter)

		mockPresenter.On("ErrorResponse", mock.AnythingOfType("*gin.Context"), errs.ErrInvalidInput).Return()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orders/process", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessOrders(c)

		mockPresenter.AssertExpectations(t)
	})

	t.Run("Empty orders array", func(t *testing.T) {
		mockProcessor := new(MockOrderProcessor)
		mockPresenter := new(MockPresenter)

		handler := handler.NewOrderHandler(mockProcessor, mockPresenter)

		mockPresenter.On("ErrorResponse", mock.AnythingOfType("*gin.Context"), errs.ErrInvalidInput).Return()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody, _ := json.Marshal([]*model.InputOrder{})
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orders/process", bytes.NewBuffer(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessOrders(c)

		mockPresenter.AssertExpectations(t)
	})

	t.Run("Invalid input order - negative quantity", func(t *testing.T) {
		mockProcessor := new(MockOrderProcessor)
		mockPresenter := new(MockPresenter)

		handler := handler.NewOrderHandler(mockProcessor, mockPresenter)

		inputData := []*model.InputOrder{
			{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               -1,
				UnitPrice:         50.0,
				TotalPrice:        100.0,
			},
		}

		mockPresenter.On("ErrorResponse", mock.AnythingOfType("*gin.Context"), mock.Anything).Return()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody, _ := json.Marshal(inputData)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orders/process", bytes.NewBuffer(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessOrders(c)

		mockPresenter.AssertExpectations(t)
	})

	t.Run("Order processor returns error", func(t *testing.T) {
		mockProcessor := new(MockOrderProcessor)
		mockPresenter := new(MockPresenter)

		handler := handler.NewOrderHandler(mockProcessor, mockPresenter)

		inputData := []*model.InputOrder{
			{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         50.0,
				TotalPrice:        100.0,
			},
		}

		processingError := errors.New("processing failed")

		mockProcessor.On("ProcessOrders", mock.AnythingOfType("[]*entity.InputOrder")).Return(([]*entity.CleanedOrder)(nil), processingError)
		mockPresenter.On("ErrorResponse", mock.AnythingOfType("*gin.Context"), processingError).Return()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody, _ := json.Marshal(inputData)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orders/process", bytes.NewBuffer(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessOrders(c)

		mockProcessor.AssertExpectations(t)
		mockPresenter.AssertExpectations(t)
	})

	t.Run("Invalid price values", func(t *testing.T) {
		mockProcessor := new(MockOrderProcessor)
		mockPresenter := new(MockPresenter)

		handler := handler.NewOrderHandler(mockProcessor, mockPresenter)

		inputData := []*model.InputOrder{
			{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         -10.0,
				TotalPrice:        100.0,
			},
		}

		mockPresenter.On("ErrorResponse", mock.AnythingOfType("*gin.Context"), mock.Anything).Return()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		requestBody, _ := json.Marshal(inputData)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orders/process", bytes.NewBuffer(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessOrders(c)

		mockPresenter.AssertExpectations(t)
	})
}

func TestOrderHandler_ProcessOrders_ValidationErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validationTests := []struct {
		name      string
		inputData []*model.InputOrder
		wantError bool
	}{
		{
			name: "Valid order",
			inputData: []*model.InputOrder{
				{
					No:                1,
					PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
					Qty:               2,
					UnitPrice:         50.0,
					TotalPrice:        100.0,
				},
			},
			wantError: false,
		},
		{
			name: "Invalid order - zero quantity",
			inputData: []*model.InputOrder{
				{
					No:                1,
					PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
					Qty:               0,
					UnitPrice:         50.0,
					TotalPrice:        100.0,
				},
			},
			wantError: true,
		},
		{
			name: "Invalid order - empty platform product id",
			inputData: []*model.InputOrder{
				{
					No:                1,
					PlatformProductId: "",
					Qty:               2,
					UnitPrice:         50.0,
					TotalPrice:        100.0,
				},
			},
			wantError: true,
		},
		{
			name: "Invalid order - negative unit price",
			inputData: []*model.InputOrder{
				{
					No:                1,
					PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
					Qty:               2,
					UnitPrice:         -50.0,
					TotalPrice:        100.0,
				},
			},
			wantError: true,
		},
		{
			name: "Invalid order - negative total price",
			inputData: []*model.InputOrder{
				{
					No:                1,
					PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
					Qty:               2,
					UnitPrice:         50.0,
					TotalPrice:        -100.0,
				},
			},
			wantError: true,
		},
	}

	for _, tt := range validationTests {
		t.Run(tt.name, func(t *testing.T) {
			mockProcessor := new(MockOrderProcessor)
			mockPresenter := new(MockPresenter)

			handler := handler.NewOrderHandler(mockProcessor, mockPresenter)

			if tt.wantError {
				mockPresenter.On("ErrorResponse", mock.AnythingOfType("*gin.Context"), mock.Anything).Return()
			} else {
				mockProcessor.On("ProcessOrders", mock.AnythingOfType("[]*entity.InputOrder")).Return([]*entity.CleanedOrder{}, nil)
				mockPresenter.On("SuccessResponse", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("[]*model.CleanedOrder")).Return()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			requestBody, _ := json.Marshal(tt.inputData)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orders/process", bytes.NewBuffer(requestBody))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.ProcessOrders(c)

			mockPresenter.AssertExpectations(t)
			if !tt.wantError {
				mockProcessor.AssertExpectations(t)
			}
		})
	}
}

func BenchmarkOrderHandler_ProcessOrders(b *testing.B) {
	gin.SetMode(gin.TestMode)

	mockProcessor := new(MockOrderProcessor)
	mockPresenter := new(MockPresenter)

	handler := handler.NewOrderHandler(mockProcessor, mockPresenter)

	inputData := []*model.InputOrder{
		{
			No:                1,
			PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
			Qty:               2,
			UnitPrice:         50.0,
			TotalPrice:        100.0,
		},
	}

	expectedResult := []*entity.CleanedOrder{
		{
			No:         1,
			ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
			MaterialId: "FG0A-CLEAR",
			ModelId:    "IPHONE16PROMAX",
			Qty:        2,
			UnitPrice:  value_object.MustNewPrice(50.0),
			TotalPrice: value_object.MustNewPrice(100.0),
		},
	}

	mockProcessor.On("ProcessOrders", mock.AnythingOfType("[]*entity.InputOrder")).Return(expectedResult, nil)
	mockPresenter.On("SuccessResponse", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("[]*model.CleanedOrder")).Return()

	requestBody, _ := json.Marshal(inputData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orders/process", bytes.NewBuffer(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ProcessOrders(c)
	}
}
