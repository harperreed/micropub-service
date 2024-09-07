package git

import (
	"os"
	"path/filepath"
	"testing"
	"fmt"
	"strings"
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
	mockGitOps := NewMockGitOperations()
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
		{
			name: "Missing title",
			content: map[string]interface{}{
				"properties": map[string]interface{}{
					"content": []interface{}{"This is a test post without a title"},
				},
			},
			wantErr: false,
		},
		{
			name: "Empty properties",
			content: map[string]interface{}{
				"properties": map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name:    "Nil content",
			content: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mockGitOps.CreatePost(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePost() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				// Check if post was created in mock storage
				url, ok := tt.content["url"].(string)
				if !ok {
					t.Errorf("CreatePost() did not set URL in content map")
				} else {
					if _, exists := mockGitOps.Posts[url]; !exists {
						t.Errorf("CreatePost() post not created: %s", url)
					}
				}
			}
		})
	}
}

// TestUpdatePost tests the UpdatePost function
func TestUpdatePost(t *testing.T) {
	mockGitOps := NewMockGitOperations()

	// Create a test post
	initialPost := map[string]interface{}{
		"properties": map[string]interface{}{
			"title":   []interface{}{"Initial Title"},
			"content": []interface{}{"Initial content"},
		},
	}
	err := mockGitOps.CreatePost(initialPost)
	if err != nil {
		t.Fatalf("Failed to create initial post: %v", err)
	}

	initialURL := initialPost["url"].(string)

	tests := []struct {
		name    string
		content map[string]interface{}
		wantErr bool
		check   func(*testing.T, *MockGitOperations, string)
	}{
		{
			name: "Update title and content",
			content: map[string]interface{}{
				"url": initialURL,
				"properties": map[string]interface{}{
					"title":   []interface{}{"Updated Title"},
					"content": []interface{}{"Updated content"},
				},
			},
			wantErr: false,
			check: func(t *testing.T, m *MockGitOperations, url string) {
				content := m.Posts[url]
				if !strings.Contains(content, "title: Updated Title") {
					t.Errorf("Title not updated correctly")
				}
				if !strings.Contains(content, "Updated content") {
					t.Errorf("Content not updated correctly")
				}
			},
		},
		{
			name: "Update title only",
			content: map[string]interface{}{
				"url": initialURL,
				"properties": map[string]interface{}{
					"title": []interface{}{"New Title"},
				},
			},
			wantErr: false,
			check: func(t *testing.T, m *MockGitOperations, url string) {
				content := m.Posts[url]
				if !strings.Contains(content, "title: New Title") {
					t.Errorf("Title not updated correctly")
				}
				if !strings.Contains(content, "Updated content") {
					t.Errorf("Content should not have changed")
				}
			},
		},
		{
			name: "Update content only",
			content: map[string]interface{}{
				"url": initialURL,
				"properties": map[string]interface{}{
					"content": []interface{}{"New content"},
				},
			},
			wantErr: false,
			check: func(t *testing.T, m *MockGitOperations, url string) {
				content := m.Posts[url]
				if !strings.Contains(content, "title: New Title") {
					t.Errorf("Title should not have changed")
				}
				if !strings.Contains(content, "New content") {
					t.Errorf("Content not updated correctly")
				}
			},
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
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mockGitOps.UpdatePost(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && tt.check != nil {
				tt.check(t, mockGitOps, initialURL)
			}
		})
	}
}

// TestDeletePost tests the DeletePost function
func TestDeletePost(t *testing.T) {
	mockGitOps := NewMockGitOperations()

	// Create test posts
	testPosts := []map[string]interface{}{
		{
			"properties": map[string]interface{}{
				"title":   []interface{}{"Test Post 1"},
				"content": []interface{}{"Content of test post 1"},
			},
		},
		{
			"properties": map[string]interface{}{
				"title":   []interface{}{"Test Post 2"},
				"content": []interface{}{"Content of test post 2"},
			},
		},
	}

	for _, post := range testPosts {
		err := mockGitOps.CreatePost(post)
		if err != nil {
			t.Fatalf("Failed to create test post: %v", err)
		}
	}

	tests := []struct {
		name    string
		content map[string]interface{}
		wantErr bool
		check   func(*testing.T, *MockGitOperations)
	}{
		{
			name: "Delete existing post",
			content: map[string]interface{}{
				"url": testPosts[0]["url"],
			},
			wantErr: false,
			check: func(t *testing.T, m *MockGitOperations) {
				if _, exists := m.Posts[testPosts[0]["url"].(string)]; exists {
					t.Errorf("Post was not deleted")
				}
			},
		},
		{
			name: "Delete non-existent post",
			content: map[string]interface{}{
				"url": "/non-existent.md",
			},
			wantErr: true,
			check:   nil,
		},
		{
			name: "Invalid input - missing URL",
			content: map[string]interface{}{
				"properties": map[string]interface{}{},
			},
			wantErr: true,
			check:   nil,
		},
		{
			name:    "Invalid input - nil content",
			content: nil,
			wantErr: true,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mockGitOps.DeletePost(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeletePost() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && tt.check != nil {
				tt.check(t, mockGitOps)
			}
		})
	}

	// Check if the remaining post is still there
	if _, exists := mockGitOps.Posts[testPosts[1]["url"].(string)]; !exists {
		t.Errorf("Undeleted post was accidentally removed")
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
		{"Filename with mixed case", "MiXeD CaSe FiLe", "mixed-case-file"},
		{"Filename with multiple spaces", "Multiple   Spaces   Here", "multiple-spaces-here"},
		{"Filename with leading and trailing hyphens", "-test-file-", "test-file"},
		{"Filename with underscores", "test_file_name", "test_file_name"},
		{"Filename with numbers", "file123name456", "file123name456"},
		{"Filename with non-ASCII characters", "テスト_ファイル.txt", "txt"},
		{"Empty filename", "", ""},
		{"Filename with only spaces", "   ", ""},
		{"Very long filename", strings.Repeat("a", 100), strings.Repeat("a", 100)},
		{"Filename with multiple consecutive hyphens", "test---file", "test-file"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeFilename(tt.filename)
			if got != tt.want {
				t.Errorf("sanitizeFilename() = %v, want %v", got, tt.want)
			}
			if len(got) > 0 {
				if got[0] == '-' || got[len(got)-1] == '-' {
					t.Errorf("sanitizeFilename() returned string with leading or trailing hyphen: %v", got)
				}
				if strings.Contains(got, "--") {
					t.Errorf("sanitizeFilename() returned string with consecutive hyphens: %v", got)
				}
			}
		})
	}
}

// Add more helper functions for testing as needed
