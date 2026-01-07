package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginationResponse represents a paginated response
type PaginationResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Message string      `json:"message,omitempty"`
}

// SendSuccess sends a successful response
func SendSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// SendSuccessWithMessage sends a successful response with a message
func SendSuccessWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// SendCreated sends a created response
func SendCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
		Message: "Created successfully",
	})
}

// SendPaginated sends a paginated response
func SendPaginated(c *gin.Context, data interface{}, total, page, limit int) {
	c.JSON(http.StatusOK, PaginationResponse{
		Success: true,
		Data:    data,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}

// SendError sends an error response
func SendError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   err.Error(),
	})
}

// SendErrorMessage sends an error response with a custom message
func SendErrorMessage(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
	})
}

// SendBadRequest sends a bad request response
func SendBadRequest(c *gin.Context, err error) {
	SendError(c, http.StatusBadRequest, err)
}

// SendUnauthorized sends an unauthorized response
func SendUnauthorized(c *gin.Context, message string) {
	SendErrorMessage(c, http.StatusUnauthorized, message)
}

// SendForbidden sends a forbidden response
func SendForbidden(c *gin.Context, message string) {
	SendErrorMessage(c, http.StatusForbidden, message)
}

// SendNotFound sends a not found response
func SendNotFound(c *gin.Context, message string) {
	SendErrorMessage(c, http.StatusNotFound, message)
}

// SendInternalServerError sends an internal server error response
func SendInternalServerError(c *gin.Context, err error) {
	SendError(c, http.StatusInternalServerError, err)
}
