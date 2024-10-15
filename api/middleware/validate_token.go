package middleware

import (
	"ems/infrastructure/config"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func ValidateToken(tokenValue string) (uint, error) {

	token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.Config.JwtSecretKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	userID, ok := claims["userID"].(float64)

	if !ok || userID == 0 {
		return 0, errors.New("invalid token")
	}

	return uint(userID), nil
}

func GetUserClaims(c *gin.Context) (*UserMiddleWareClaims, error) {
	userClaims, isSuccess := c.Get("user")

	if !isSuccess {
		return nil, errors.New("user fetching error in middleware")
	}

	user, ok := userClaims.(*UserMiddleWareClaims)

	if !ok {
		return nil, errors.New("user type conversion error")
	}
	return user, nil
}
