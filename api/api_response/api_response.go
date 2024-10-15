package api_response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, ApiResponse{
		Message: message,
		Data:    data,
	})
}

func BadRequestError(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, ApiResponse{
		Message: message,
	})
}

func InternalServerError(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, ApiResponse{
		Message: message,
	})
}

func UnauthorizedError(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, ApiResponse{
		Message: message,
	})
}
