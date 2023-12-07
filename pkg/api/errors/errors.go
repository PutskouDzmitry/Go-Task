package errors

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"task/internal/delivery/errors"

	e "github.com/pkg/errors"
)

type errorResponse struct {
	Code    int    `json:"status_code"`
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, err error) {
	var response errorResponse
	switch {
	case e.Is(err, errors.ErrBadInput):
		response = errorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("bad request - %s", err),
		}
	case e.Is(err, errors.ErrAccessDenied):
		fallthrough
	case e.Is(err, errors.ErrNotAllowed):
		fallthrough
	case e.Is(err, errors.ErrWrongData):
		response = errorResponse{
			Code:    http.StatusForbidden,
			Message: fmt.Sprintf("forbidden - %s", err),
		}
	case e.Is(err, errors.ErrNotFound):
		response = errorResponse{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("not found - %s", err),
		}
	case e.Is(err, errors.ErrInternal):
		fallthrough
	default:
		response = errorResponse{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("internal server error - %s", err),
		}
	}

	c.AbortWithStatusJSON(response.Code, response)
}

func AbortWithErrorResponse(c *gin.Context, r errorResponse) {
	c.AbortWithStatusJSON(r.Code, r)
}

func AbortWithBadRequest(c *gin.Context, err error) {
	AbortWithErrorResponse(c, errorResponse{
		Code:    http.StatusBadRequest,
		Message: "bad request: " + err.Error(),
	})
}

func AbortWithUnauthorized(c *gin.Context, err error) {
	AbortWithErrorResponse(c, errorResponse{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized: " + err.Error(),
	})
}

func NewForbiddenError(c *gin.Context) {
	AbortWithErrorResponse(c, errorResponse{
		Message: "forbidden: not enough rights",
		Code:    http.StatusForbidden,
	})
}
