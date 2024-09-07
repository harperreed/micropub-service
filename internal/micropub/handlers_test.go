package micropub

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	// "os"
	// "path/filepath"
	"strings"
	"testing"
	// "io/ioutil"

	"github.com/harperreed/micropub-service/internal/git"
	"github.com/labstack/echo/v5"
)

type MockGitOperations struct {
    CreatePostError error
    UpdatePostError error
    DeletePostError error
    MockFileContent string
}

// Add this struct
type MockEventEmitter struct {
	EmitCalled bool
}

func (m *MockEventEmitter) Emit(event interface{}) {
	m.EmitCalled = true
}

func (m *MockGitOperations) CreatePost(content map[string]interface{}) error {
	if m.CreatePostError != nil {
		return m.CreatePostError
	}
	// Simulate the behavior of the real CreatePost function
	properties, ok := content["properties"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid properties")
	}
	contentValue, ok := properties["content"].([]interface{})
	if !ok || len(contentValue) == 0 {
		return fmt.Errorf("invalid content")
	}
	// Simulate setting the URL
	content["url"] = "https://example.com/new-post"
	return nil
}

func (m *MockGitOperations) UpdatePost(content map[string]interface{}) error {
    if m.UpdatePostError != nil {
        return m.UpdatePostError
    }

    // We're not using filepath in our mock, so we can remove this check
    _, ok := content["url"].(string)
    if !ok {
        return fmt.Errorf("invalid URL")
    }

    properties, ok := content["properties"].(map[string]interface{})
    if !ok {
        return fmt.Errorf("invalid properties data")
    }

    // Instead of reading from a file, use the mock content
    existingContent := m.MockFileContent
    if existingContent == "" {
        existingContent = "---\ntitle: Initial Title\n---\nInitial content"
    }

    // Parse the existing frontmatter
    frontmatter, _, err := SplitFrontmatterAndContent(existingContent)
    if err != nil {
        return err
    }

    // Update the frontmatter with new values
    for key, value := range properties {
        if key != "content" {
            if slice, ok := value.([]interface{}); ok && len(slice) > 0 {
                frontmatter[key] = fmt.Sprintf("%v", slice[0])
            } else {
                frontmatter[key] = fmt.Sprintf("%v", value)
            }
        }
    }

    // Get the updated content
    var updatedContent string
    if content, ok := properties["content"]; ok {
        if slice, ok := content.([]interface{}); ok && len(slice) > 0 {
            updatedContent = fmt.Sprintf("%v", slice[0])
        } else {
            updatedContent = fmt.Sprintf("%v", content)
        }
    }

    // Create the updated content with frontmatter
    fullContent := CreateContentWithFrontmatter(frontmatter, updatedContent)

    // Instead of writing to a file, update the mock content
    m.MockFileContent = fullContent

    return nil
}


func (m *MockGitOperations) DeletePost(content map[string]interface{}) error {
	// Mock implementation
	if m.DeletePostError != nil {
		return m.DeletePostError
	}
	return nil
}

func TestHandleMicropubUpdate(t *testing.T) {
    e := echo.New()

    req := httptest.NewRequest(http.MethodPut, "/micropub", strings.NewReader(`
        {
            "action": "update",
            "url": "https://example.com/2023-05-01-test-post.md",
            "replace": {
                "content": ["This is an updated test post"],
                "title": ["Updated Title"]
            }
        }`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    // Mock git operations
    originalGitOps := git.GitOps
    mockGitOps := &MockGitOperations{
        MockFileContent: "---\ntitle: Initial Title\n---\nInitial content",
    }
    git.GitOps = mockGitOps
    defer func() { git.GitOps = originalGitOps }()

    if err := HandleMicropubUpdate(c); err != nil {
        t.Fatalf("HandleMicropubUpdate failed: %v", err)
    }

    if rec.Code != http.StatusOK {
        t.Errorf("Expected status OK; got %v", rec.Code)
    }

    expected := "Post updated successfully"
    if strings.TrimSpace(rec.Body.String()) != expected {
        t.Errorf("Expected body %q; got %q", expected, rec.Body.String())
    }

    // Check if the content was updated correctly
    expectedContent := "---\ntitle: Updated Title\n---\nThis is an updated test post"
    if mockGitOps.MockFileContent != expectedContent {
        t.Errorf("File content not updated correctly. Expected:\n%s\nGot:\n%s", expectedContent, mockGitOps.MockFileContent)
    }
}

func SplitFrontmatterAndContent(content string) (map[string]string, string, error) {
    parts := strings.SplitN(content, "---", 3)
    if len(parts) != 3 {
        return nil, "", fmt.Errorf("invalid content format")
    }

    frontmatter := make(map[string]string)
    for _, line := range strings.Split(strings.TrimSpace(parts[1]), "\n") {
        kv := strings.SplitN(line, ":", 2)
        if len(kv) == 2 {
            frontmatter[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
        }
    }

    return frontmatter, strings.TrimSpace(parts[2]), nil
}

func CreateContentWithFrontmatter(frontmatter map[string]string, content string) string {
    var sb strings.Builder
    sb.WriteString("---\n")
    for k, v := range frontmatter {
        sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
    }
    sb.WriteString("---\n")
    sb.WriteString(content)
    return sb.String()
}

func TestHandleMicropubCreate(t *testing.T) {
    // Set up test directory
    testDir := t.TempDir()
    originalRepoPath := git.RepoPath
    git.RepoPath = testDir
    defer func() { git.RepoPath = originalRepoPath }()

    // Use mock git operations
    originalGitOps := git.GitOps
    git.GitOps = &git.MockGitOperations{}
    defer func() { git.GitOps = originalGitOps }()
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

	t.Run("SuccessfulCreateWithoutURL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/micropub", strings.NewReader(`{"type":["h-entry"],"properties":{"content":["Ahoy, world!"],"category":["test","micropub"]}}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockEmitter := &MockEventEmitter{}
		SetEventEmitter(mockEmitter)
		defer SetEventEmitter(nil)

		originalGitOps := git.GitOps
		git.GitOps = &MockGitOperations{}
		defer func() { git.GitOps = originalGitOps }()

		if err := HandleMicropubCreate(c); err != nil {
			t.Fatalf("HandleMicropubCreate failed: %v", err)
		}

		if !mockEmitter.EmitCalled {
			t.Errorf("Expected event emitter to be called")
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status Created; got %v", rec.Code)
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

	t.Run("SuccessfulCreateFormEncoded", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("h", "entry")
		formData.Set("content", "Ahoy, world!")
		formData.Add("category[]", "test")
		formData.Add("category[]", "micropub")

		req := httptest.NewRequest(http.MethodPost, "/micropub", strings.NewReader(formData.Encode()))
		req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Mock git operations
		originalGitOps := git.GitOps
		git.GitOps = &MockGitOperations{}
		defer func() { git.GitOps = originalGitOps }()

		if err := HandleMicropubCreate(c); err != nil {
			t.Fatalf("HandleMicropubCreate failed: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status Created; got %v", rec.Code)
		}
	})

	t.Run("MissingType", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/micropub", strings.NewReader(`{"properties":{"content":["Ahoy, world!"]}}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := HandleMicropubCreate(c)

		if err == nil {
			t.Errorf("Expected error for missing type, got nil")
		}
		if httperr, ok := err.(*echo.HTTPError); ok {
			if httperr.Code != http.StatusBadRequest {
				t.Errorf("Expected status BadRequest; got %v", httperr.Code)
			}
		} else {
			t.Errorf("Expected *echo.HTTPError, got %T", err)
		}
	})

	t.Run("MissingContent", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/micropub", strings.NewReader(`{"type":["h-entry"],"properties":{}}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := HandleMicropubCreate(c)

		if err == nil {
			t.Errorf("Expected error for missing content, got nil")
		}
		if httperr, ok := err.(*echo.HTTPError); ok {
			if httperr.Code != http.StatusBadRequest {
				t.Errorf("Expected status BadRequest; got %v", httperr.Code)
			}
		} else {
			t.Errorf("Expected *echo.HTTPError, got %T", err)
		}
	})

	t.Run("UnsupportedContentType", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/micropub", strings.NewReader(`{"type":["h-entry"],"properties":{"content":["Ahoy, world!"]}}`))
		req.Header.Set(echo.HeaderContentType, "application/xml")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := HandleMicropubCreate(c)

		if err == nil {
			t.Errorf("Expected error for unsupported content type, got nil")
		}
		if httperr, ok := err.(*echo.HTTPError); ok {
			if httperr.Code != http.StatusUnsupportedMediaType {
				t.Errorf("Expected status UnsupportedMediaType; got %v", httperr.Code)
			}
		} else {
			t.Errorf("Expected *echo.HTTPError, got %T", err)
		}
	})

	t.Run("SuccessfulCreateWithEventEmitter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/micropub", strings.NewReader(`{"type":["h-entry"],"properties":{"content":["Ahoy, world!"],"category":["test","micropub"]}}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockEmitter := &MockEventEmitter{}
		SetEventEmitter(mockEmitter)
		defer SetEventEmitter(nil)

		originalGitOps := git.GitOps
		git.GitOps = &MockGitOperations{}
		defer func() { git.GitOps = originalGitOps }()

		if err := HandleMicropubCreate(c); err != nil {
			t.Fatalf("HandleMicropubCreate failed: %v", err)
		}

		if !mockEmitter.EmitCalled {
			t.Errorf("Expected event emitter to be called")
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status Created; got %v", rec.Code)
		}
	})

	t.Run("SuccessfulCreateWithEventEmitterNoURL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/micropub", strings.NewReader(`{"type":["h-entry"],"properties":{"content":["Ahoy, world!"],"category":["test","micropub"]}}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockEmitter := &MockEventEmitter{}
		SetEventEmitter(mockEmitter)
		defer SetEventEmitter(nil)

		if err := HandleMicropubCreate(c); err != nil {
			t.Fatalf("HandleMicropubCreate failed: %v", err)
		}

		if !mockEmitter.EmitCalled {
			t.Errorf("Expected event emitter to be called")
		}
	})

	t.Run("CreatePostError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/micropub", strings.NewReader(`{"type":["h-entry"],"properties":{"content":["Ahoy, world!"],"category":["test","micropub"]}}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		originalGitOps := git.GitOps
		git.GitOps = &MockGitOperations{CreatePostError: errors.New("failed to create post")}
		defer func() { git.GitOps = originalGitOps }()

		err := HandleMicropubCreate(c)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if httperr, ok := err.(*echo.HTTPError); ok {
			if httperr.Code != http.StatusInternalServerError {
				t.Errorf("Expected status InternalServerError; got %v", httperr.Code)
			}
		} else {
			t.Errorf("Expected *echo.HTTPError, got %T", err)
		}
	})
}

func TestHandleMicropubUpdateInvalid(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/micropub", strings.NewReader(`{"invalid json`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := HandleMicropubUpdate(c)

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

func TestHandleMicropubUpdateScenarios(t *testing.T) {
    e := echo.New()

    // Set up test directory
    testDir := t.TempDir()
    originalRepoPath := git.RepoPath
    git.RepoPath = testDir
    defer func() { git.RepoPath = originalRepoPath }()

    // Use mock git operations
    originalGitOps := git.GitOps
    git.GitOps = &git.MockGitOperations{}
    defer func() { git.GitOps = originalGitOps }()

    t.Run("SuccessfulUpdate", func(t *testing.T) {
           mockGitOps := &MockGitOperations{
               MockFileContent: "---\ntitle: Initial Title\n---\nInitial content",
           }
           git.GitOps = mockGitOps
           defer func() { git.GitOps = originalGitOps }()

           req := httptest.NewRequest(http.MethodPut, "/micropub", strings.NewReader(`
               {
                   "action": "update",
                   "url": "https://example.com/2023-05-01-test-post.md",
                   "replace": {
                       "content": ["Updated content"],
                       "title": ["Updated Title"]
                   }
               }`))
           req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
           rec := httptest.NewRecorder()
           c := e.NewContext(req, rec)

           if err := HandleMicropubUpdate(c); err != nil {
               t.Fatalf("HandleMicropubUpdate failed: %v", err)
           }

           if rec.Code != http.StatusOK {
               t.Errorf("Expected status OK; got %v", rec.Code)
           }

           expectedContent := "---\ntitle: Updated Title\n---\nUpdated content"
           if mockGitOps.MockFileContent != expectedContent {
               t.Errorf("File content not updated correctly. Expected:\n%s\nGot:\n%s", expectedContent, mockGitOps.MockFileContent)
           }
       })

	t.Run("UpdatePostError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/micropub", strings.NewReader(`{"action":"update","url":"https://example.com/post1","replace":{"content":["Updated content"]}}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		originalGitOps := git.GitOps
		git.GitOps = &MockGitOperations{UpdatePostError: errors.New("failed to update post")}
		defer func() { git.GitOps = originalGitOps }()

		err := HandleMicropubUpdate(c)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if httperr, ok := err.(*echo.HTTPError); ok {
			if httperr.Code != http.StatusInternalServerError {
				t.Errorf("Expected status InternalServerError; got %v", httperr.Code)
			}
		} else {
			t.Errorf("Expected *echo.HTTPError, got %T", err)
		}
	})
}

// Add this function
func TestSetEventEmitter(t *testing.T) {
	mockEmitter := &MockEventEmitter{}
	SetEventEmitter(mockEmitter)
	if eventEmitter != mockEmitter {
		t.Errorf("Expected eventEmitter to be set to mockEmitter")
	}
}
