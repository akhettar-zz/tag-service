package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

const (
	Referer   = "Referer"
	UserAgent = "User-Agent"
)

// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
// By default gin.DefaultWriter = os.Stdout.
func Logger() gin.HandlerFunc {
	return LogWithWriter(gin.DefaultWriter)
}

// This
func LogWithWriter(out io.Writer, notlogged ...string) gin.HandlerFunc {

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			end := time.Now()
			latency := end.Sub(start)

			clientIP := c.ClientIP()

			referer := c.GetHeader(Referer)
			userAgent := c.GetHeader(UserAgent)
			method := c.Request.Method
			statusCode := c.Writer.Status()
			comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

			if raw != "" {
				path = path + "?" + raw
			}

			fmt.Fprintf(out, "[GIN] %v |%3d | %13v | X-Forwarded-For: %15s| Referer: %s | User-Agent: %s | %-7s %s\n%s",
				end.Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				clientIP, referer, userAgent,
				method,
				path,
				comment,
			)
		}
	}
}
