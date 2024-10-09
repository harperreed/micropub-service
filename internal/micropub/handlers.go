package micropub

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

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
	if err != nil {
		return err
	}

	// Check if required fields are present
	if content["type"] == nil || len(content["type"].([]interface{})) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing 'type' field")
	}

	properties, ok := content["properties"].(map[string]interface{})
	if !ok || properties["content"] == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing or invalid 'content' field")
	}

	// Handle multiple media attachments
	if photos, ok := properties["photo"]; ok {
		if photoSlice, ok := photos.([]interface{}); ok {
			properties["photo"] = photoSlice
		} else if photoStr, ok := photos.(string); ok {
			properties["photo"] = []string{photoStr}
		}
	}

	err = git.GitOps.CreatePost(content)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create post: "+err.Error())
	}

	if eventEmitter != nil {
		postID, _ := content["url"].(string)
		if postID == "" {
			postID = "unknown"
		}
		eventEmitter.Emit(PostEvent{Type: "create", PostID: postID})
	}

	return c.String(http.StatusCreated, "Post created successfully")
}

func HandleMicropubUpdate(c echo.Context) error {
	content, err := parseContent(c)
	if err != nil {
		return err
	}

	if content["action"] != "update" || content["url"] == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid update request")
	}

	// Handle 'replace' and 'add' actions
	properties := make(map[string]interface{})
	if replace, ok := content["replace"].(map[string]interface{}); ok {
		for key, value := range replace {
			if sliceValue, ok := value.([]interface{}); ok && len(sliceValue) > 0 {
				properties[key] = sliceValue[0]
			} else {
				properties[key] = value
			}
		}
	}
	if add, ok := content["add"].(map[string]interface{}); ok {
		for key, value := range add {
			if existingValue, exists := properties[key]; exists {
				if existingSlice, ok := existingValue.([]interface{}); ok {
					properties[key] = append(existingSlice, value)
				} else {
					properties[key] = []interface{}{existingValue, value}
				}
			} else {
				properties[key] = value
			}
		}
	}
	content["properties"] = properties

	err = git.GitOps.UpdatePost(content)
	if err != nil {
		if errors.Is(err, git.ErrPostNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Post not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update post: "+err.Error())
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
		if errors.Is(err, git.ErrPostNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Post not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete post: "+err.Error())
	}

	return c.String(http.StatusOK, "Post deleted successfully")
}

func parseContent(c echo.Context) (map[string]interface{}, error) {
	req := c.Request()
	contentType := req.Header.Get("Content-Type")
	var content map[string]interface{}

	switch {
	case strings.HasPrefix(contentType, "multipart/form-data"):
		if err := req.ParseMultipartForm(32 << 20); err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Error parsing multipart form: "+err.Error())
		}
		content = make(map[string]interface{})
		properties := make(map[string]interface{})
		for key, values := range req.MultipartForm.Value {
			if key == "h" {
				content["type"] = []interface{}{fmt.Sprintf("h-%s", values[0])}
			} else if strings.HasSuffix(key, "[]") {
				properties[strings.TrimSuffix(key, "[]")] = values
			} else if len(values) == 1 {
				properties[key] = values[0]
			} else {
				properties[key] = values
			}
		}
		for key, fileHeaders := range req.MultipartForm.File {
			files := make([]string, len(fileHeaders))
			for i, fileHeader := range fileHeaders {
				file, err := fileHeader.Open()
				if err != nil {
					return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error opening uploaded file: "+err.Error())
				}
				defer file.Close()
				// Here you would typically save the file and store its path or URL
				files[i] = fileHeader.Filename // For demonstration, we're just storing the filename
			}
			properties[key] = files
		}
		content["properties"] = properties
	case contentType == "application/x-www-form-urlencoded":
		if err := req.ParseForm(); err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Error parsing form data: "+err.Error())
		}
		content = make(map[string]interface{})
		properties := make(map[string]interface{})
		for key, values := range req.Form {
			if key == "h" {
				content["type"] = []interface{}{fmt.Sprintf("h-%s", values[0])}
			} else if strings.HasSuffix(key, "[]") {
				properties[strings.TrimSuffix(key, "[]")] = values
			} else if len(values) == 1 {
				properties[key] = values[0]
			} else {
				properties[key] = values
			}
		}
		content["properties"] = properties
	case contentType == "application/json":
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
