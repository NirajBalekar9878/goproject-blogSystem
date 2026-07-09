package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs HTTP requests in the format:
// METHOD /path -> STATUS StatusText (Duration)
// Example: POST /blogs -> 201 Created (20ms)
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		statusText := http.StatusText(statusCode)

		// Format duration e.g., 20ms or 2.5ms
		durationStr := fmt.Sprintf("%dms", duration.Milliseconds())
		if duration.Milliseconds() == 0 {
			durationStr = fmt.Sprintf("%dµs", duration.Microseconds())
		}

		log.Printf("%s %s -> %d %s (%s)",
			c.Request.Method,
			c.Request.URL.Path,
			statusCode,
			statusText,
			durationStr,
		)
	}
}
