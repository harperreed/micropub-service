package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/patrickmn/go-cache"
	"github.com/pocketbase/pocketbase/models"
)

var userRoleCache *cache.Cache

func init() {
	userRoleCache = cache.New(5*time.Minute, 10*time.Minute)
}

func RoleAuthorization(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, _ := c.Get("user").(*models.Record)

			if user == nil {
				return c.String(http.StatusUnauthorized, "You must be logged in to access this resource")
			}

			userRole := getUserRole(user.Id)
			for _, role := range allowedRoles {
				if userRole == role {
					return next(c)
				}
			}

			return c.String(http.StatusForbidden, "You are not authorized to access this resource")
		}
	}
}

func getUserRole(userId string) string {
	if cachedRole, found := userRoleCache.Get(userId); found {
		return cachedRole.(string)
	}

	// If not found in cache, fetch from database and cache it
	// This is a placeholder - replace with actual database query
	role := "user" // Default role
	userRoleCache.Set(userId, role, cache.DefaultExpiration)
	return role
}
