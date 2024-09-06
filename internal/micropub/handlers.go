package micropub

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v5"
	"github.com/harperreed/micropub-service/internal/git"
)

// FileEvent represents a file-related event
type FileEvent struct {
	Type     string // e.g., "upload", "delete"
	Filename string
}

// EventEmitter is an interface for emitting events
type EventEmitter interface {
	Emit(event interface{})
}

var eventEmitter EventEmitter

func HandleMicropubCreate(c echo.Context) error {
	content, err := parseContent(c)
	if err != nil {
		return err
	}

	err = git.GitOps.CreatePost(content)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create post")
	}

	// Emit file upload event
	if eventEmitter != nil {
		filename := content["filename"].(string)
		eventEmitter.Emit(FileEvent{Type: "upload", Filename: filename})
	}

	return c.String(http.StatusCreated, "Post created successfully")
}

func HandleMicropubUpdate(c echo.Context) error {
	content, err := parseContent(c)
	if err != nil {
		return err
	}

	err = git.GitOps.UpdatePost(content)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to update post")
	}

	return c.String(http.StatusOK, "Post updated successfully")
}

func HandleMicropubDelete(c echo.Context) error {
	content, err := parseContent(c)
	if err != nil {
		return err
	}

	err = git.GitOps.DeletePost(content)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete post")
	}

	// Emit file delete event
	if eventEmitter != nil {
		filename := content["url"].(string) // Assuming the URL is the filename
		eventEmitter.Emit(FileEvent{Type: "delete", Filename: filename})
	}

	return c.String(http.StatusOK, "Post deleted successfully")
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

// SetEventEmitter sets the event emitter for the package
func SetEventEmitter(emitter EventEmitter) {
	eventEmitter = emitter
}
