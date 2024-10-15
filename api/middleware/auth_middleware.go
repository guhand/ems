package middleware

import (
	"ems/domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	userRepository domain.UserRepository
}

func NewMiddleware(userRepository domain.UserRepository) *Middleware {
	return &Middleware{userRepository}
}

type UserMiddleWareClaims struct {
	ID                 uint
	RoleID             uint
	DepartmentID       *uint
	DepartmentMemberID *uint
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token missing"})
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		userID, err := ValidateToken(token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		user, err := m.userRepository.GetUserByID(userID)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if user == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if user.Token == nil || *user.Token == "" || *user.Token != token {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		userClaims := &UserMiddleWareClaims{
			ID:                 user.ID,
			RoleID:             user.RoleID,
			DepartmentID:       user.DepartmentID,
			DepartmentMemberID: user.DepartmentMemberID,
		}

		c.Set("user", userClaims)

		c.Next()
	}
}
