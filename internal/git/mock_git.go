package git

import (
	"fmt"
	"time"
	"path/filepath"
	"os"
	"log"
)

// MockGitOperations is a mock implementation of the GitOperations interface for testing purposes
type MockGitOperations struct {
	Posts map[string]string // map to store posts with URL as key and content as value
}

// NewMockGitOperations creates a new instance of MockGitOperations
func NewMockGitOperations() *MockGitOperations {
	return &MockGitOperations{
		Posts: make(map[string]string),
	}
}

// CreatePost simulates creating a new post
func (m *MockGitOperations) CreatePost(content map[string]interface{}) error {
	properties, ok := content["properties"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid properties")
	}

	title := "Untitled Post"
	if titleValue, ok := properties["title"]; ok {
		if titleSlice, ok := titleValue.([]interface{}); ok && len(titleSlice) > 0 {
			title, _ = titleSlice[0].(string)
		}
	}

	body := ""
	if contentValue, ok := properties["content"]; ok {
		if contentSlice, ok := contentValue.([]interface{}); ok && len(contentSlice) > 0 {
			body, _ = contentSlice[0].(string)
		}
	}

	if body == "" {
		return fmt.Errorf("missing content")
	}

	filename := fmt.Sprintf("%s-%s.md", time.Now().Format("2006-01-02"), sanitizeFilename(title))
	url := fmt.Sprintf("/%s", filename)

	m.Posts[url] = fmt.Sprintf("---\ntitle: %s\ndate: %s\n---\n\n%s", title, time.Now().Format(time.RFC3339), body)
	content["url"] = url

	log.Printf("Created post: %s", url)
	return nil
}

// UpdatePost simulates updating an existing post
func (m *MockGitOperations) UpdatePost(content map[string]interface{}) error {
	url, ok := content["url"].(string)
	if !ok {
		return fmt.Errorf("invalid URL")
	}

	properties, ok := content["properties"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid properties")
	}

	existingContent, exists := m.Posts[url]
	if !exists {
		return fmt.Errorf("post not found")
	}

	frontmatter, oldContent, err := SplitFrontmatterAndContent(existingContent)
	if err != nil {
		return err
	}

	if title, ok := properties["title"].([]interface{}); ok && len(title) > 0 {
		frontmatter["title"] = title[0]
	}

	if newContent, ok := properties["content"].([]interface{}); ok && len(newContent) > 0 {
		oldContent = newContent[0].(string)
	}

	m.Posts[url] = CreateContentWithFrontmatter(frontmatter, oldContent)

	log.Printf("Updated post: %s", url)
	return nil
}

// DeletePost simulates deleting a post
func (m *MockGitOperations) DeletePost(content map[string]interface{}) error {
	url, ok := content["url"].(string)
	if !ok {
		return fmt.Errorf("invalid URL")
	}

	if _, exists := m.Posts[url]; !exists {
		return fmt.Errorf("post not found")
	}

	delete(m.Posts, url)

	log.Printf("Deleted post: %s", url)
	return nil
}

// InitializeRepo simulates initializing a git repository
func (m *MockGitOperations) InitializeRepo() error {
	log.Println("Simulating repository initialization")
	return nil
}

// gitAdd simulates adding a file to git
func (m *MockGitOperations) gitAdd(filename string) error {
	log.Printf("Simulating git add: %s", filename)
	return nil
}

// gitCommit simulates committing changes to git
func (m *MockGitOperations) gitCommit(message string) error {
	log.Printf("Simulating git commit: %s", message)
	return nil
}

// gitPush simulates pushing changes to a remote repository
func (m *MockGitOperations) gitPush() error {
	log.Println("Simulating git push")
	return nil
}
