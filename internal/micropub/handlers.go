package micropub

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v5"
	"log"
)

func HandleMicropubCreate(c echo.Context) error {
	content, err := parseContent(c)
	if err != nil {
		return err
	}

	// TODO: Implement CreatePost functionality
	log.Printf("Attempting to create post with content: %v", content)
	return c.String(http.StatusNotImplemented, "Create functionality not yet implemented")
}

func HandleMicropubUpdate(c echo.Context) error {
	content, err := parseContent(c)
	if err != nil {
		return err
	}

	// TODO: Implement UpdatePost functionality
	log.Printf("Attempting to update post with content: %v", content)
	return c.String(http.StatusNotImplemented, "Update functionality not yet implemented")
}

func HandleMicropubDelete(c echo.Context) error {
	content, err := parseContent(c)
	if err != nil {
		return err
	}

	// TODO: Implement DeletePost functionality
	log.Printf("Attempting to delete post with content: %v", content)
	return c.String(http.StatusNotImplemented, "Delete functionality not yet implemented")
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
