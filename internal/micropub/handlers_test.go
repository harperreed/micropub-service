package micropub

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/harperreed/micropub-service/internal/git"
)

func TestCreatePost(t *testing.T) {
	// ... (keep existing TestCreatePost function)
}

func TestUpdatePost(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-repo")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set the repoPath to the temporary directory
	originalRepoPath := git.RepoPath
	git.RepoPath = tempDir
	defer func() { git.RepoPath = originalRepoPath }()

	// Initialize a test Git repository
	if err := git.InitializeRepo(); err != nil {
		t.Fatalf("Failed to initialize test repository: %v", err)
	}

	// Create a test post
	createTestPost(t)

	// Test for JSON request
	t.Run("JSON", func(t *testing.T) {
		jsonBody := map[string]interface{}{
			"action":  "update",
			"url":     "https://example.com/2023-05-01-test-post.md",
			"content": "This is an updated test post",
		}
		jsonBytes, _ := json.Marshal(jsonBody)

		req, err := http.NewRequest("PUT", "/micropub", bytes.NewBuffer(jsonBytes))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := echo.HandlerFunc(HandleMicropubUpdate)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Post updated successfully and pushed to Git repository"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}

func TestDeletePost(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-repo")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set the repoPath to the temporary directory
	originalRepoPath := git.RepoPath
	git.RepoPath = tempDir
	defer func() { git.RepoPath = originalRepoPath }()

	// Initialize a test Git repository
	if err := git.InitializeRepo(); err != nil {
		t.Fatalf("Failed to initialize test repository: %v", err)
	}

	// Create a test post
	createTestPost(t)

	// Test for JSON request
	t.Run("JSON", func(t *testing.T) {
		jsonBody := map[string]interface{}{
			"action": "delete",
			"url":    "https://example.com/2023-05-01-test-post.md",
		}
		jsonBytes, _ := json.Marshal(jsonBody)

		req, err := http.NewRequest("DELETE", "/micropub", bytes.NewBuffer(jsonBytes))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := echo.HandlerFunc(HandleMicropubDelete)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := "Post deleted successfully and pushed to Git repository"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}

func createTestPost(t *testing.T) {
	content := map[string]interface{}{
		"title":   "Test Post",
		"content": "This is a test post",
	}
	// TODO: Implement CreatePost functionality
	// For now, we'll just log the content
	log.Printf("Creating test post with content: %v", content)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}
}
