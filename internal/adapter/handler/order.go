package handler

import (
	"order-placement-system/internal/adapter/handler/model"
	"order-placement-system/internal/adapter/presenter"
	usecase "order-placement-system/internal/usecases/interfaces"
	"order-placement-system/pkg/log"

	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	orderProcessor usecase.OrderProcessorUseCase
	presenter      presenter.OrderPresenter
}

type OrderHandlerInterface interface {
	ProcessOrders(c *gin.Context)
}

func NewOrderHandler(
	orderProcessor usecase.OrderProcessorUseCase,
	presenter presenter.OrderPresenter,
) OrderHandlerInterface {
	return &orderHandler{
		orderProcessor: orderProcessor,
		presenter:      presenter,
	}
}
func (h *orderHandler) ProcessOrders(c *gin.Context) {

	var inputOrderModels []*model.InputOrder
	req, err := new(model.InputOrder).Parse(c)
	if err != nil {
		log.Errorf("failed to parse request body", log.E(err))
		h.presenter.ErrorResponse(c, err)
		return
	}

	inputOrderModels = req
	result, err := h.orderProcessor.ProcessOrders(model.ToEntity(inputOrderModels))
	if err != nil {
		log.Errorf("failed to process orders", log.E(err))
		h.presenter.ErrorResponse(c, err)
		return
	}

	h.presenter.SuccessResponse(c, model.FromEntities(result))
}
