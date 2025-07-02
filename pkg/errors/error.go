package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrNotFound            = errors.New("entity not found")
	ErrAlreadyExists       = errors.New("entity already exists")
	ErrInvalidInput        = errors.New("invalid input")
	ErrUnauthorized        = errors.New("unauthorized access")
	ErrInternalServer      = errors.New("internal server error")
	ErrConflict            = errors.New("conflict")
	ErrForbidden           = errors.New("forbidden")
	ErrBadRequest          = errors.New("bad request")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	ErrTooManyRequests     = errors.New("too many requests")
)

func MapJsonError(c *gin.Context, err error) {
	switch err {
	case ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case ErrInvalidInput:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case ErrAlreadyExists:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case ErrUnprocessableEntity:
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	case ErrUnauthorized:
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case ErrForbidden:
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case ErrConflict:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case ErrTooManyRequests:
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}
}
