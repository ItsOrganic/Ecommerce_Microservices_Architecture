package middleware

import (
	"strconv"
	"time"
	"user-service/metrics" // Assuming this is where your custom metrics are defined

	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware is a Gin middleware for collecting Prometheus metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record the start time
		start := time.Now()

		// Process the request
		c.Next()

		// Calculate the latency (duration)
		duration := time.Since(start).Seconds()

		// Extract the request path and method
		path := c.FullPath() // Use c.FullPath() for routes
		if path == "" {
			path = "undefined"
		}
		method := c.Request.Method
		statusCode := c.Writer.Status()
		status := strconv.Itoa(statusCode)

		// Update Prometheus metrics
		metrics.HttpRequestCounter.WithLabelValues(path, method, status).Inc()      // Increment request count
		metrics.HttpRequestDuration.WithLabelValues(path, method).Observe(duration) // Record request duration
	}
}
