package presenter

import (
	"net/http"
	"order-placement-system/pkg/errors"

	"github.com/gin-gonic/gin"
)

type OrderPresenter interface {
	SuccessResponse(c *gin.Context, data interface{})
	ErrorResponse(c *gin.Context, err error)
}

type orderPresenter struct{}

func NewOrderPresenter() OrderPresenter {
	return &orderPresenter{}
}

func (p *orderPresenter) SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}

func (p *orderPresenter) ErrorResponse(c *gin.Context, err error) {
	errors.MapJsonError(c, err)
}
