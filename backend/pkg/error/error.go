package error

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Error constructors
func BadRequest(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func InternalError(message string, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
		Err:        err,
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	}
}

// HandleError logs and responds with appropriate error
func HandleError(c *gin.Context, err *AppError) {
	// Log internal errors with details
	if err.StatusCode >= 500 {
		if err.Err != nil {
			log.Printf("[ERROR] %s: %v", err.Message, err.Err)
		} else {
			log.Printf("[ERROR] %s", err.Message)
		}
	}

	// Return generic message to client
	c.JSON(err.StatusCode, gin.H{
		"error": err.Message,
	})
}

// ErrorHandler middleware for panic recovery
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
