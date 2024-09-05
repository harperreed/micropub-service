package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/micropub", handleMicropubCreate)
		e.Router.PUT("/micropub", handleMicropubUpdate)
		e.Router.DELETE("/micropub", handleMicropubDelete)
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func handleMicropubCreate(c *core.ServeEvent) error {
	content, err := parseContent(c)
	if err != nil {
		return err
	}

	err = createPost(content)
	if err != nil {
		return c.Router.Context().String(http.StatusInternalServerError, fmt.Sprintf("Error creating post: %v", err))
	}

	return c.Router.Context().String(http.StatusCreated, "Post created successfully and pushed to Git repository")
}

func handleMicropubUpdate(c *core.ServeEvent) error {
	content, err := parseContent(c)
	if err != nil {
		return err
	}

	err = updatePost(content)
	if err != nil {
		return c.Router.Context().String(http.StatusInternalServerError, fmt.Sprintf("Error updating post: %v", err))
	}

	return c.Router.Context().String(http.StatusOK, "Post updated successfully and pushed to Git repository")
}

func handleMicropubDelete(c *core.ServeEvent) error {
	content, err := parseContent(c)
	if err != nil {
		return err
	}

	err = deletePost(content)
	if err != nil {
		return c.Router.Context().String(http.StatusInternalServerError, fmt.Sprintf("Error deleting post: %v", err))
	}

	return c.Router.Context().String(http.StatusOK, "Post deleted successfully and pushed to Git repository")
}

func parseContent(c *core.ServeEvent) (map[string]interface{}, error) {
	req := c.Router.Context().Request()
	contentType := req.Header.Get("Content-Type")
	var content map[string]interface{}

	switch contentType {
	case "application/x-www-form-urlencoded":
		if err := req.ParseForm(); err != nil {
			return nil, c.Router.Context().String(http.StatusBadRequest, "Error parsing form data")
		}
		content = parseFormToMap(req.PostForm)
	case "application/json":
		if err := c.Router.Context().Bind(&content); err != nil {
			return nil, c.Router.Context().String(http.StatusBadRequest, "Error parsing JSON")
		}
	default:
		return nil, c.Router.Context().String(http.StatusUnsupportedMediaType, "Unsupported Content-Type")
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
