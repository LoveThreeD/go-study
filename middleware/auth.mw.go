package middleware

import (
	"github.com/gin-gonic/gin"
	"sTest/pkg/auth"
	"sTest/pkg/response"
)

// UserAuthMiddleware Valid user token (jwt)
func UserAuthMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		// XXX 该处获取Token可以从redis中获取
		// get token
		token := c.GetHeader("Authorization")
		if token == "" {
			response.ResultFailed(c, response.ErrNoPermission)
			c.Next()
			return
		}

		// parse token
		claims, err := auth.ParseToken(token)
		if err != nil {
			response.ResultFailed(c, response.ErrNoPermission)
			c.Next()
			return
		}

		// set context
		c.Set("token", claims)
		c.Next()
	}
}
