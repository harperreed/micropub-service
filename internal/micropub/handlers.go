package micropub

import (
	"encoding/json"
	"net/http"
	"net/url"
	"fmt"
	"github.com/harperreed/micropub-service/internal/git"
	"github.com/labstack/echo/v5"
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

type PostEvent struct {
	Type   string
	PostID string
}

var eventEmitter EventEmitter

func HandleMicropubCreate(c echo.Context) error {
	content, err := parseContent(c)
	// if err != nil {
	//     return echo.NewHTTPError(http.StatusBadRequest, "Invalid request: "+err.Error())
	// }

    if err != nil {
        return err
    }

    if content["type"] == nil || len(content["type"].([]interface{})) == 0 {
        return echo.NewHTTPError(http.StatusBadRequest, "Missing 'type' field")
    }

    properties, ok := content["properties"].(map[string]interface{})
    if !ok || properties["content"] == nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Missing or invalid 'content' field")
    }


	if eventEmitter != nil {
		postID := content["url"].(string) // Assuming the URL is set after creation
		eventEmitter.Emit(PostEvent{Type: "create", PostID: postID})
	}

	err = git.GitOps.CreatePost(content)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create post: "+err.Error())
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

	if _, ok := content["url"]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing URL for delete action")
	}

	err = git.GitOps.DeletePost(content)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete post")
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
            return nil, echo.NewHTTPError(http.StatusBadRequest, "Error parsing form data: "+err.Error())
        }
        content = make(map[string]interface{})
        for key, values := range req.Form {
            if len(values) == 1 {
                content[key] = values[0]
            } else {
                content[key] = values
            }
        }
        // Special handling for 'h' field in form data
        if h, ok := content["h"]; ok {
            content["type"] = []interface{}{fmt.Sprintf("h-%s", h)}
        }
        // Move content to properties
        if _, ok := content["content"]; ok {
            properties := make(map[string]interface{})
            for k, v := range content {
                if k != "h" && k != "type" {
                    properties[k] = v
                }
            }
            content["properties"] = properties
            delete(content, "content")
        }
	case "application/json":
		if err := json.NewDecoder(req.Body).Decode(&content); err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Error parsing JSON: "+err.Error())
		}
	default:
		return nil, echo.NewHTTPError(http.StatusUnsupportedMediaType, "Unsupported Content-Type")
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
