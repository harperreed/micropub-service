package micropub

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"fmt"
	"github.com/labstack/echo/v5"
	"github.com/harperreed/micropub-service/internal/git"
)

type MockGitOperations struct{}

func (m *MockGitOperations) CreatePost(content map[string]interface{}) error {
    // Check if the content is valid
    if content["type"] == nil || len(content["type"].([]interface{})) == 0 {
        return fmt.Errorf("missing 'type' field")
    }
    properties, ok := content["properties"].(map[string]interface{})
    if !ok || properties["content"] == nil {
        return fmt.Errorf("missing or invalid 'content' field")
    }
    return nil
}

func (m *MockGitOperations) UpdatePost(content map[string]interface{}) error {
    // Mock implementation
    return nil
}

func (m *MockGitOperations) DeletePost(content map[string]interface{}) error {
    // Mock implementation
    return nil
}

func TestHandleMicropubUpdate(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new request
	req := httptest.NewRequest(http.MethodPut, "/micropub", strings.NewReader(`{"action":"update","url":"https://example.com/2023-05-01-test-post.md","content":"This is an updated test post"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Create a ResponseRecorder to record the response
	rec := httptest.NewRecorder()

	// Create a new Echo context
	c := e.NewContext(req, rec)

	// Mock git operations
	originalGitOps := git.GitOps
	git.GitOps = &MockGitOperations{}
	defer func() { git.GitOps = originalGitOps }()

	// Call the handler
	if err := HandleMicropubUpdate(c); err != nil {
		t.Fatalf("HandleMicropubUpdate failed: %v", err)
	}

	// Check the status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", rec.Code)
	}

	// Check the response body
	expected := "Post updated successfully"
	if strings.TrimSpace(rec.Body.String()) != expected {
		t.Errorf("Expected body %q; got %q", expected, rec.Body.String())
	}
}

func TestHandleMicropubCreate(t *testing.T) {
    // Create a new Echo instance
    e := echo.New()

    // Test case 1: Successful post creation
    t.Run("SuccessfulCreate", func(t *testing.T) {
        // Create a new request
        req := httptest.NewRequest(http.MethodPost, "/micropub", strings.NewReader(`{"type":["h-entry"],"properties":{"content":["Ahoy, world!"],"category":["test","micropub"]}}`))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

        // Create a ResponseRecorder to record the response
        rec := httptest.NewRecorder()

        // Create a new Echo context
        c := e.NewContext(req, rec)

        // Mock git operations
        originalGitOps := git.GitOps
        git.GitOps = &MockGitOperations{}
        defer func() { git.GitOps = originalGitOps }()

        // Call the handler
        if err := HandleMicropubCreate(c); err != nil {
            t.Fatalf("HandleMicropubCreate failed: %v", err)
        }

        // Check the status code
        if rec.Code != http.StatusCreated {
            t.Errorf("Expected status Created; got %v", rec.Code)
        }

        // Check the response body
        expected := "Post created successfully"
        if strings.TrimSpace(rec.Body.String()) != expected {
            t.Errorf("Expected body %q; got %q", expected, rec.Body.String())
        }
    })


    // Test case 2: Invalid JSON
    t.Run("InvalidJSON", func(t *testing.T) {
        req := httptest.NewRequest(http.MethodPost, "/micropub", strings.NewReader(`{"invalid json`))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)

        err := HandleMicropubCreate(c)

        if err == nil {
            t.Errorf("Expected error for invalid JSON, got nil")
        }
        if httperr, ok := err.(*echo.HTTPError); ok {
            if httperr.Code != http.StatusBadRequest {
                t.Errorf("Expected status BadRequest; got %v", httperr.Code)
            }
        } else {
            t.Errorf("Expected *echo.HTTPError, got %T", err)
        }
    })
}

func TestHandleMicropubDelete(t *testing.T) {
    e := echo.New()

    // Test case 1: Successful post deletion
    t.Run("SuccessfulDelete", func(t *testing.T) {
        req := httptest.NewRequest(http.MethodDelete, "/micropub", strings.NewReader(`{"action":"delete","url":"https://example.com/2023-05-01-test-post.md"}`))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)

        originalGitOps := git.GitOps
        git.GitOps = &MockGitOperations{}
        defer func() { git.GitOps = originalGitOps }()

        if err := HandleMicropubDelete(c); err != nil {
            t.Fatalf("HandleMicropubDelete failed: %v", err)
        }

        if rec.Code != http.StatusOK {
            t.Errorf("Expected status OK; got %v", rec.Code)
        }

        expected := "Post deleted successfully"
        if strings.TrimSpace(rec.Body.String()) != expected {
            t.Errorf("Expected body %q; got %q", expected, rec.Body.String())
        }
    })

    // Test case 2: Missing URL
    t.Run("MissingURL", func(t *testing.T) {
            req := httptest.NewRequest(http.MethodDelete, "/micropub", strings.NewReader(`{"action":"delete"}`))
            req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
            rec := httptest.NewRecorder()
            c := e.NewContext(req, rec)

            err := HandleMicropubDelete(c)

            if err == nil {
                t.Errorf("Expected error for missing URL, got nil")
            }
            if httperr, ok := err.(*echo.HTTPError); ok {
                if httperr.Code != http.StatusBadRequest {
                    t.Errorf("Expected status BadRequest; got %v", httperr.Code)
                }
            } else {
                t.Errorf("Expected *echo.HTTPError, got %T", err)
            }
        })
}
