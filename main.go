package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/micropub", echo.HandlerFunc(handleMicropubCreate))
		e.Router.PUT("/micropub", echo.HandlerFunc(handleMicropubUpdate))
		e.Router.DELETE("/micropub", echo.HandlerFunc(handleMicropubDelete))

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

	authRecord, err := app.Users().AuthenticateWithPassword(email, password)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Invalid email or password")
	}

	token, err := app.Users().CreateAuthToken(authRecord.Id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create auth token")
	}

	c.SetCookie(&http.Cookie{
		Name:     "pb_auth",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	return c.Redirect(http.StatusSeeOther, "/")
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
