package auth

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type MockAuthMiddleware struct {
	ValidToken string
	ValidScope string
}

func NewMockAuthMiddleware(validToken, validScope string) *MockAuthMiddleware {
	return &MockAuthMiddleware{
		ValidToken: validToken,
		ValidScope: validScope,
	}
}

func (m *MockAuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token != m.ValidToken {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid access token")
		}

		scope := c.QueryParam("scope")
		if scope != m.ValidScope {
			return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
		}

		return next(c)
	}
}
