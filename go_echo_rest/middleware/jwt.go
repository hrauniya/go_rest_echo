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

var AccessJWTSecret = []byte(getAccessTokenJWTSecret())
var RefreshJWTSecret = []byte(getRefreshTokenJWTSecret())


type Claims struct {
	UserID uint `json:"user_id"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID uint   `json:"user_id"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func getAccessTokenJWTSecret() string {
	secret := os.Getenv("secret")
	return secret
}

func getRefreshTokenJWTSecret() string {
	refreshtokensecret := os.Getenv("refreshsecretkey")
	return refreshtokensecret
}

func GenerateAccessToken(user model.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(AccessJWTSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshToken(user model.User) (string, error) {
	claims := &RefreshClaims{
		UserID: user.ID,
		Type:   "refresh", 
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), 
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(RefreshJWTSecret)
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
			if claims.Type != "access" {
				return nil, errors.New("invalid token type for this endpoint")
			}
			return AccessJWTSecret, nil
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


func ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		
		if claims.Type != "refresh" {
			return nil, errors.New("invalid token type: expected refresh token")
		}
		return 	RefreshJWTSecret, nil
	})

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired refresh token")
	}

	if !token.Valid {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid refresh token")
	}

	return claims, nil
}








