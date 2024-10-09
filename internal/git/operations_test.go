package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var origRepoPath string
var testRepoPath string

func setupTestEnvironment() {
	origRepoPath = RepoPath
	testRepoPath = filepath.Join(os.TempDir(), "test-repo")
	err := os.MkdirAll(testRepoPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create test directory: %v\n", err)
		os.Exit(1)
	}
	RepoPath = testRepoPath
	GitOps = &MockGitOperations{}
}

func teardownTestEnvironment() {
	os.RemoveAll(testRepoPath)
	RepoPath = origRepoPath
}

func TestMain(m *testing.M) {
	setupTestEnvironment()
	code := m.Run()
	teardownTestEnvironment()
	os.Exit(code)
}

func createTestFile(name, content string) error {
	return ioutil.WriteFile(filepath.Join(testRepoPath, name), []byte(content), 0644)
}

func readTestFile(name string) (string, error) {
	content, err := ioutil.ReadFile(filepath.Join(testRepoPath, name))
	if err != nil {
		return "", err
	}
	return string(content), nil
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

func TestHandleMergeConflicts(t *testing.T) {
	setupTestEnvironment()
	defer teardownTestEnvironment()

	// Create a file with initial content
	err := createTestFile("conflict.txt", "Initial content")
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Simulate a merge conflict
	conflictContent := `<<<<<<< HEAD
Modified content in the current branch
=======
Modified content in the other branch
>>>>>>> other-branch
`
	err = createTestFile("conflict.txt", conflictContent)
	if err != nil {
		t.Fatalf("Failed to create conflict in test file: %v", err)
	}

	// Test handling merge conflicts
	err = GitOps.HandleMergeConflicts()
	if err != nil {
		t.Errorf("HandleMergeConflicts() error = %v", err)
	}

	// Verify that the conflict has been resolved
	content, err := readTestFile("conflict.txt")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}
	if strings.Contains(content, "<<<<<<<") || strings.Contains(content, "=======") || strings.Contains(content, ">>>>>>>") {
		t.Errorf("Merge conflict was not resolved properly")
	}
}

func TestBranchOperations(t *testing.T) {
	setupTestEnvironment()
	defer teardownTestEnvironment()

	// Test creating a new branch
	err := GitOps.CreateBranch("new-feature")
	if err != nil {
		t.Errorf("CreateBranch() error = %v", err)
	}

	// Test switching to the new branch
	err = GitOps.SwitchBranch("new-feature")
	if err != nil {
		t.Errorf("SwitchBranch() error = %v", err)
	}

	// Create a file in the new branch
	err = createTestFile("feature.txt", "New feature content")
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Switch back to the main branch
	err = GitOps.SwitchBranch("main")
	if err != nil {
		t.Errorf("SwitchBranch() error = %v", err)
	}

	// Verify that the file doesn't exist in the main branch
	_, err = readTestFile("feature.txt")
	if err == nil {
		t.Errorf("File should not exist in the main branch")
	}
}

func TestLargeFileUploads(t *testing.T) {
	setupTestEnvironment()
	defer teardownTestEnvironment()

	// Create a large file (100MB)
	largeContent := make([]byte, 100*1024*1024) // 100MB of zero bytes
	err := ioutil.WriteFile(filepath.Join(testRepoPath, "large_file.bin"), largeContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create large test file: %v", err)
	}

	// Test adding and committing the large file
	err = GitOps.AddFile("large_file.bin")
	if err == nil {
		t.Errorf("AddFile() should return an error for large files")
	}

	// Verify that the error message indicates a file size limit
	if !strings.Contains(err.Error(), "file size limit") {
		t.Errorf("Expected error message to mention file size limit, got: %v", err)
	}
}

func TestGitPushErrors(t *testing.T) {
	setupTestEnvironment()
	defer teardownTestEnvironment()

	// Test network error during push
	GitOps.(*MockGitOperations).SetPushBehavior(errors.New("network error"))
	err := GitOps.Push()
	if err == nil || !strings.Contains(err.Error(), "network error") {
		t.Errorf("Push() should return a network error, got: %v", err)
	}

	// Test authentication failure during push
	GitOps.(*MockGitOperations).SetPushBehavior(errors.New("authentication failed"))
	err = GitOps.Push()
	if err == nil || !strings.Contains(err.Error(), "authentication failed") {
		t.Errorf("Push() should return an authentication error, got: %v", err)
	}

	// Test successful push
	GitOps.(*MockGitOperations).SetPushBehavior(nil)
	err = GitOps.Push()
	if err != nil {
		t.Errorf("Push() should succeed, got error: %v", err)
	}
}

// Add more helper functions for testing as needed
