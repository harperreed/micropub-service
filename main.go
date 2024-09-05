package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/patrickmn/go-cache"
)

var userRoleCache *cache.Cache

func init() {
	userRoleCache = cache.New(5*time.Minute, 10*time.Minute)
}

func roleAuthorization(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			app := c.Get("app").(*pocketbase.PocketBase)
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

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/micropub", echo.HandlerFunc(handleMicropubCreate), roleAuthorization("admin", "editor"))
		e.Router.PUT("/micropub", echo.HandlerFunc(handleMicropubUpdate), roleAuthorization("admin", "editor"))
		e.Router.DELETE("/micropub", echo.HandlerFunc(handleMicropubDelete), roleAuthorization("admin"))

		// Add routes for login
		e.Router.GET("/login", echo.HandlerFunc(handleLoginPage))
		e.Router.POST("/login", echo.HandlerFunc(handleLogin))

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func handleLoginPage(c echo.Context) error {
	return c.File("login.html")
}

func handleLogin(c echo.Context) error {
	app := c.Get("app").(*pocketbase.PocketBase)

	email := c.FormValue("email")
	password := c.FormValue("password")

	authRecord, err := app.Dao().FindAuthRecordByEmail("users", email)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Invalid email or password")
	}

	if !authRecord.ValidatePassword(password) {
		return c.String(http.StatusUnauthorized, "Invalid email or password")
	}

	token, err := app.NewAuthToken(authRecord)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create auth token")
	}

	// Cache the user's role
	role := authRecord.Get("role").(string) // Assuming the role is stored in a "role" field
	userRoleCache.Set(authRecord.Id, role, cache.DefaultExpiration)

	c.SetCookie(&http.Cookie{
		Name:     "pb_auth",
		Value:    token.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
		"role":  role,
	})
}

func createPost(content map[string]interface{}) error {
	// Implement the createPost logic here
	return nil
}

func updatePost(content map[string]interface{}) error {
	// Implement the updatePost logic here
	return nil
}

func deletePost(content map[string]interface{}) error {
	// Implement the deletePost logic here
	return nil
}

func handleMicropubCreate(c echo.Context) error {
	content, err := parseContent(c)
	if err != nil {
		return err
	}

	err = createPost(content)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error creating post: %v", err))
	}

	return c.String(http.StatusCreated, "Post created successfully and pushed to Git repository")
}

func handleMicropubUpdate(c echo.Context) error {
	content, err := parseContent(c)
	if err != nil {
		return err
	}

	err = updatePost(content)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error updating post: %v", err))
	}

	return c.String(http.StatusOK, "Post updated successfully and pushed to Git repository")
}

func handleMicropubDelete(c echo.Context) error {
	content, err := parseContent(c)
	if err != nil {
		return err
	}

	err = deletePost(content)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error deleting post: %v", err))
	}

	return c.String(http.StatusOK, "Post deleted successfully and pushed to Git repository")
}

func parseContent(c echo.Context) (map[string]interface{}, error) {
	req := c.Request()
	contentType := req.Header.Get("Content-Type")
	var content map[string]interface{}

	switch contentType {
	case "application/x-www-form-urlencoded":
		if err := req.ParseForm(); err != nil {
			return nil, c.String(http.StatusBadRequest, "Error parsing form data")
		}
		content = parseFormToMap(req.PostForm)
	case "application/json":
		if err := json.NewDecoder(req.Body).Decode(&content); err != nil {
			return nil, c.String(http.StatusBadRequest, "Error parsing JSON")
		}
	default:
		return nil, c.String(http.StatusUnsupportedMediaType, "Unsupported Content-Type")
	}

	return content, nil
}

func parseFormToMap(form url.Values) map[string]interface{} {
	result := make(map[string]interface{})
	for key, values := range form {
		if len(values) == 1 {
			result[key] = values[0]
		} else {
			result[key] = values
		}
	}
	return result
}
