package micropub

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/harperreed/micropub-service/internal/git"
)

// MockGitOperations is a mock implementation of git operations
type MockGitOperations struct{}

func (m *MockGitOperations) UpdatePost(content map[string]interface{}) error {
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
