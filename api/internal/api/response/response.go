package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   struct {
		Code    int    `json:"code"`
		Details string `json:"details"`
	} `json:"error"`
}

func JSONSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, SuccessResponse{
		Status: "success",
		Data:   data,
	})
}

func JSONCreated(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Status: "success",
		Data:   data,
	})
}

func JSONError(c *gin.Context, code int, message string, details string) {
	c.JSON(code, ErrorResponse{
		Status:  "error",
		Message: message,
		Error: struct {
			Code    int    `json:"code"`
			Details string `json:"details"`
		}{
			Code:    code,
			Details: details,
		},
	})
}
