package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/patrickmn/go-cache"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"

	"github.com/harperreed/micropub-service/internal/config"
	"github.com/harperreed/micropub-service/internal/events"
	"github.com/harperreed/micropub-service/internal/git"
	"github.com/harperreed/micropub-service/internal/micropub"
)

var userRoleCache *cache.Cache

func init() {
	userRoleCache = cache.New(5*time.Minute, 10*time.Minute)
}

func setupFileCleanup(emitter *events.EventEmitter) {
	emitter.On("file", func(event interface{}) {
		fileEvent := event.(events.FileEvent)
		log.Printf("File event received: %v", fileEvent)
		// Implement your file cleanup logic here
	})
}

func roleAuthorization(allowedRoles ...string) echo.MiddlewareFunc {
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

func handleLoginPage(c echo.Context) error {
	return c.File("templates/login.html")
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

	token, err := tokens.NewRecordAuthToken(app, authRecord)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create auth token")
	}

	// Cache the user's role
	role := authRecord.Get("role").(string) // Assuming the role is stored in a "role" field
	userRoleCache.Set(authRecord.Id, role, cache.DefaultExpiration)

	c.SetCookie(&http.Cookie{
		Name:     "pb_auth",
		Value:    token,
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
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Use the configuration
	log.Printf("Git repository path: %s", cfg.GitRepoPath)

	app := pocketbase.New()

	// Initialize event emitter
	eventEmitter := events.NewEventEmitter()
	micropub.SetEventEmitter(eventEmitter)

	// Initialize Git repository
	if err := git.InitializeRepo(); err != nil {
    	log.Fatalf("Failed to initialize Git repository: %v", err)
	}

	// Set up file cleanup process
	setupFileCleanup(eventEmitter)

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/micropub", echo.HandlerFunc(micropub.HandleMicropubCreate), roleAuthorization("admin", "editor"))
		e.Router.PUT("/micropub", echo.HandlerFunc(micropub.HandleMicropubUpdate), roleAuthorization("admin", "editor"))
		e.Router.DELETE("/micropub", echo.HandlerFunc(micropub.HandleMicropubDelete), roleAuthorization("admin"))

		// Add routes for login
		e.Router.GET("/login", echo.HandlerFunc(handleLoginPage))
		e.Router.POST("/login", echo.HandlerFunc(handleLogin))

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
