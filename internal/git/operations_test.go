package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup
	origRepoPath := RepoPath
	RepoPath = filepath.Join(os.TempDir(), "test-repo")
	err := os.MkdirAll(RepoPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create test directory: %v\n", err)
		os.Exit(1)
	}
	GitOps = &MockGitOperations{}

	// Run tests
	code := m.Run()

	// Teardown
	os.RemoveAll(RepoPath)
	RepoPath = origRepoPath

	os.Exit(code)
}

// TestCreatePost tests the CreatePost function
func TestCreatePost(t *testing.T) {
	tests := []struct {
		name        string
		content     map[string]interface{}
		wantErr     bool
		wantContent string
	}{
		{
			name: "Valid post with title and content",
			content: map[string]interface{}{
				"properties": map[string]interface{}{
					"title":   []interface{}{"Test Post"},
					"content": []interface{}{"This is a test post"},
				},
			},
			wantErr:     false,
			wantContent: "---\ntitle: Test Post\ndate: ",
		},
		{
			name: "Valid post with HTML content",
			content: map[string]interface{}{
				"properties": map[string]interface{}{
					"title":   []interface{}{"HTML Post"},
					"content": []interface{}{"<p>This is an <strong>HTML</strong> post</p>"},
				},
			},
			wantErr:     false,
			wantContent: "---\ntitle: HTML Post\ndate: ",
		},
		{
			name: "Valid post with Markdown content",
			content: map[string]interface{}{
				"properties": map[string]interface{}{
					"title":   []interface{}{"Markdown Post"},
					"content": []interface{}{"# Heading\n\nThis is a **Markdown** post"},
				},
			},
			wantErr:     false,
			wantContent: "---\ntitle: Markdown Post\ndate: ",
		},
		{
			name: "Missing content",
			content: map[string]interface{}{
				"properties": map[string]interface{}{
					"title": []interface{}{"Test Post"},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid properties",
			content: map[string]interface{}{
				"properties": "not a map",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for each test
			tempDir, err := os.MkdirTemp("", "test-repo-*")
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Set the RepoPath to the temporary directory
			oldRepoPath := RepoPath
			RepoPath = tempDir
			defer func() { RepoPath = oldRepoPath }()

			err = GitOps.CreatePost(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePost() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				// Check if file was created
				url, ok := tt.content["url"].(string)
				if !ok {
					t.Errorf("CreatePost() did not set URL in content map")
				} else {
					filePath := filepath.Join(RepoPath, filepath.Base(url))
					if _, err := os.Stat(filePath); os.IsNotExist(err) {
						t.Errorf("CreatePost() file not created: %s", filePath)
					} else {
						// Check file content
						content, err := os.ReadFile(filePath)
						if err != nil {
							t.Errorf("Failed to read created file: %v", err)
						} else {
							if !strings.Contains(string(content), tt.wantContent) {
								t.Errorf("File content does not match expected. Got:\n%s\nWant to contain:\n%s", string(content), tt.wantContent)
							}
						}
					}
				}
			}
		})
	}
}

// TestUpdatePost tests the UpdatePost function
func TestUpdatePost(t *testing.T) {
	// Create a test file
	testFile := filepath.Join(RepoPath, "test-post.md")
	initialContent := `---
title: Initial Title
date: 2023-05-01T12:00:00Z
---

Initial content`
	err := os.WriteFile(testFile, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name    string
		content map[string]interface{}
		wantErr bool
	}{
		{
			name: "Update title and content",
			content: map[string]interface{}{
				"url": "/test-post.md",
				"properties": map[string]interface{}{
					"title":   []interface{}{"Updated Title"},
					"content": []interface{}{"Updated content"},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid URL",
			content: map[string]interface{}{
				"url": "/non-existent.md",
				"properties": map[string]interface{}{
					"title": []interface{}{"Updated Title"},
				},
			},
			wantErr: true,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GitOps.UpdatePost(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				// Check if file was updated
				_, err := os.ReadFile(testFile)
				if err != nil {
					t.Errorf("Failed to read updated file: %v", err)
				}
				// Add assertions to check if the content was updated correctly
			}
		})
	}
}

// TestDeletePost tests the DeletePost function
func TestDeletePost(t *testing.T) {
	// Create a test file
	testFile := filepath.Join(RepoPath, "test-delete-post.md")
	err := os.WriteFile(testFile, []byte("Test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name    string
		content map[string]interface{}
		wantErr bool
	}{
		{
			name: "Delete existing post",
			content: map[string]interface{}{
				"url": "/test-delete-post.md",
			},
			wantErr: false,
		},
		{
			name: "Delete non-existent post",
			content: map[string]interface{}{
				"url": "/non-existent.md",
			},
			wantErr: true,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GitOps.DeletePost(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeletePost() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				// Check if file was deleted
				if _, err := os.Stat(testFile); !os.IsNotExist(err) {
					t.Errorf("DeletePost() file not deleted: %s", testFile)
				}
			}
		})
	}
}

// Add more test functions for other operations here

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{"Normal filename", "My Test File", "my-test-file"},
		{"Filename with special characters", "Test@File#123", "testfile123"},
		{"Filename with spaces and hyphens", "  Test - File  ", "test-file"},
		{"Filename with only invalid characters", "@#$%^&*", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitizeFilename(tt.filename); got != tt.want {
				t.Errorf("sanitizeFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Add more helper functions for testing as needed
