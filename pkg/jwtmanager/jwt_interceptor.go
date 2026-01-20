package jwtmanager

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func (j *JWTManager) JWTInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := j.extractTokenFromHeader(c)
		if err != nil {
			c.JSON(401, gin.H{
				"error": MsgTokenRequired,
			})
			c.Abort()
			return
		}

		userID, err := j.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{
				"error": MsgInvalidToken,
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func (j *JWTManager) extractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", ErrMissingAuthHeader
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", ErrInvalidAuthFormat
	}

	return authHeader[len(bearerPrefix):], nil
}

func GetCurrentUserID(c *gin.Context) (int, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, ErrMissingUserID
	}

	id, ok := userID.(int)
	if !ok {
		return 0, ErrMissingUserID
	}

	return id, nil
}
