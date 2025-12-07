package utils

import (
	"net/http"

	"github.com/Faleeeee/URL_Shortener/internal/domain"

	"github.com/gin-gonic/gin"
)

// SendSuccess sends a success response with the standard format
func SendSuccess(c *gin.Context, message string, data interface{}, meta *domain.Meta) {
	response := domain.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Error:   nil,
		Meta:    meta,
	}
	c.JSON(http.StatusOK, response)
}

// SendError sends an error response with the standard format
func SendError(c *gin.Context, status int, message string, errCode string, errDetails string) {
	response := domain.APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
		Error: &domain.ErrorDetails{
			Code:    errCode,
			Details: errDetails,
		},
		Meta: nil,
	}
	c.JSON(status, response)
}
