package git

import (
	"os"
	"path/filepath"
	"testing"
	// "time"
	"fmt"
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
        name    string
        content map[string]interface{}
        wantErr bool
    }{
        {
            name: "Valid post",
            content: map[string]interface{}{
                "properties": map[string]interface{}{
                    "title":   []interface{}{"Test Post"},
                    "content": []interface{}{"This is a test post"},
                },
            },
            wantErr: false,
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
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := GitOps.CreatePost(tt.content)
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
