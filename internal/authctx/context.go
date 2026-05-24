package authctx

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserID(c *gin.Context) (uuid.UUID, error) {

	raw, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, errors.New("user not authenticated")
	}

	userID, err := uuid.Parse(raw.(string))
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {

	return func(c *gin.Context) {

		roleRaw, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "role not found",
			})
			return
		}

		userRole := roleRaw.(string)

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(403, gin.H{
			"error": "forbidden",
		})
	}
}
