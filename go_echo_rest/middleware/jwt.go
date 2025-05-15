package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go_echo_rest/model"
)

var JWTSecret = []byte(getJWTSecret())


type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}


func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	return secret
}


func GenerateToken(user model.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return JWTSecret, nil
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
		}

		if !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Subject)

		return next(c)
	}
}








