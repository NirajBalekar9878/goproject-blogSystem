package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse defines the standard JSON payload for successful operations.
type SuccessResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse defines the standard JSON payload for errors.
type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondOK sends a 200 OK JSON response.
func RespondOK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

// RespondCreated sends a 201 Created JSON response.
func RespondCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

// RespondBadRequest sends a 400 Bad Request JSON response.
func RespondBadRequest(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: err,
	})
}

// RespondNotFound sends a 404 Not Found JSON response.
func RespondNotFound(c *gin.Context, err string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error: err,
	})
}

// RespondInternalError sends a 500 Internal Server Error JSON response.
func RespondInternalError(c *gin.Context, err string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: err,
	})
}
