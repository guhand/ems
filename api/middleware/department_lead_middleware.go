package middleware

import (
	"ems/app/model/constant"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) DepartmentLeadMiddleware() gin.HandlerFunc {
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			return
		}

		if user.RoleID != uint(constant.Admin) && user.RoleID != uint(constant.Manager) &&
			user.RoleID != uint(constant.DepartmentLead) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not a Department Lead user"})
		}

		departmentLeadClaims := &UserMiddleWareClaims{
			ID:                 user.ID,
			RoleID:             user.RoleID,
			DepartmentID:       user.DepartmentID,
			DepartmentMemberID: user.DepartmentMemberID,
		}

		c.Set("user", departmentLeadClaims)

		c.Next()
	}
}
