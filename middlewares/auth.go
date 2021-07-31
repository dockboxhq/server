package middlewares

import "github.com/gin-gonic/gin"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// reqKey := c.Request.Header.Get("X-Auth-Key")
		// reqSecret := c.Request.Header.Get("X-Auth-Secret")

		c.Next()
	}
}
