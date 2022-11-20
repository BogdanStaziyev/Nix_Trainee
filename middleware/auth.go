package middleware

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	MW "github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
	"trainee/internal/app"
	"trainee/internal/infra/http/response"
)

type AuthMiddleware interface {
	JWT(secret string) echo.MiddlewareFunc
	ValidateJWT() echo.MiddlewareFunc
}

type authMiddleware struct {
	authService app.AuthService
	r           *redis.Client
}

func NewMiddleware(as app.AuthService, red *redis.Client) AuthMiddleware {
	return authMiddleware{
		authService: as,
		r:           red,
	}
}

func (m authMiddleware) ValidateJWT() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(*app.JwtTokenClaim)

			user, err := m.authService.ValidateJWT(claims.UID, claims.ID, false)
			if err != nil {
				return response.MessageResponse(c, http.StatusUnauthorized, "Not authorized")
			}

			c.Set("currentUser", user)

			go func() {
				m.r.Expire(fmt.Sprintf("token-%d", claims.ID), time.Minute*app.LogOF)
			}()
			return next(c)
		}
	}
}

func (m authMiddleware) JWT(secret string) echo.MiddlewareFunc {
	config := MW.JWTConfig{
		ErrorHandler: func(err error) error {
			return &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "Not authorized",
			}
		},
		SigningKey: []byte(secret),
		Claims:     &app.JwtTokenClaim{},
	}
	return MW.JWTWithConfig(config)
}
