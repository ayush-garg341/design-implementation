package middleware

import (
	"time"

	"context"
	"github.com/gin-gonic/gin"
)

func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// wrap the request context with our timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()

		// replace the request's context with our timeout ctx
		c.Request = c.Request.WithContext(ctx)

		// pass control to the next handler
		c.Next()
	}
}
