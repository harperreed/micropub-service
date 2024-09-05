package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/micropub", handleMicropub)
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func handleMicropub(c echo.Context) error {
	contentType := c.Request().Header.Get("Content-Type")
	var content map[string]interface{}

	switch contentType {
	case "application/x-www-form-urlencoded":
		if err := c.Request().ParseForm(); err != nil {
			return c.String(http.StatusBadRequest, "Error parsing form data")
		}
		content = parseFormToMap(c.Request().PostForm)
	case "application/json":
		if err := c.Bind(&content); err != nil {
			return c.String(http.StatusBadRequest, "Error parsing JSON")
		}
	default:
		return c.String(http.StatusUnsupportedMediaType, "Unsupported Content-Type")
	}

	err := createPost(content)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error creating post: %v", err))
	}

	return c.String(http.StatusCreated, "Post created successfully and pushed to Git repository")
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
