package micropub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v5"
	"github.com/harperreed/micropub-service/internal/git"
)

func HandleMicropubCreate(c echo.Context) error {
	content, err := parseContent(c)
	if err != nil {
		return err
	}

	err = git.CreatePost(content)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error creating post: %v", err))
	}

	return c.String(http.StatusCreated, "Post created successfully and pushed to Git repository")
}

func HandleMicropubUpdate(c echo.Context) error {
	content, err := parseContent(c)
	if err != nil {
		return err
	}

	// TODO: Implement UpdatePost functionality
	return c.String(http.StatusNotImplemented, "Update functionality not yet implemented")
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error updating post: %v", err))
	}

	return c.String(http.StatusOK, "Post updated successfully and pushed to Git repository")
}

func HandleMicropubDelete(c echo.Context) error {
	content, err := parseContent(c)
	if err != nil {
		return err
	}

	// TODO: Implement DeletePost functionality
	return c.String(http.StatusNotImplemented, "Delete functionality not yet implemented")
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
