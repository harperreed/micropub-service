package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type MockGitOperations struct {
	DefaultGitOperations
	pushBehavior error
}

func (m *MockGitOperations) CreatePost(content map[string]interface{}) error {
	properties := content["properties"].(map[string]interface{})
	title := "Untitled Post"
	if titleValue, ok := properties["title"]; ok {
		if titleSlice, ok := titleValue.([]interface{}); ok && len(titleSlice) > 0 {
			title, _ = titleSlice[0].(string)
		}
	}
	filename := fmt.Sprintf("%s-%s.md", time.Now().Format("2006-01-02"), sanitizeFilename(title))
	filePath := filepath.Join(RepoPath, filename)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write some content to the file
	_, err = file.WriteString("Test content")
	if err != nil {
		return err
	}

	content["url"] = fmt.Sprintf("/%s", filename)
	return nil
}

func (m *MockGitOperations) HandleMergeConflicts() error {
	// Simulate handling merge conflicts
	return nil
}

func (m *MockGitOperations) CreateBranch(name string) error {
	// Simulate creating a new branch
	return nil
}

func (m *MockGitOperations) SwitchBranch(name string) error {
	// Simulate switching to a branch
	return nil
}

func (m *MockGitOperations) AddFile(filename string) error {
	// Simulate adding a file, with a size check
	info, err := os.Stat(filepath.Join(RepoPath, filename))
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}
	if info.Size() > 50*1024*1024 { // 50MB limit
		return errors.New("file size limit exceeded")
	}
	return nil
}

func (m *MockGitOperations) Push() error {
	// Return the configured push behavior
	return m.pushBehavior
}

func (m *MockGitOperations) SetPushBehavior(err error) {
	m.pushBehavior = err
}

func (m *MockGitOperations) UpdatePost(content map[string]interface{}) error {
	url, ok := content["url"].(string)
	if !ok {
		return fmt.Errorf("invalid URL")
	}
	filePath := filepath.Join(RepoPath, filepath.Base(url))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found")
	}
	return nil
}

func (m *MockGitOperations) DeletePost(content map[string]interface{}) error {
	url, ok := content["url"].(string)
	if !ok {
		return fmt.Errorf("invalid URL")
	}
	filePath := filepath.Join(RepoPath, filepath.Base(url))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found")
	}
	return os.Remove(filePath)
}

func (m *MockGitOperations) InitializeRepo() error {
	// Simulate initializing a repo
	return nil
}

func (m *MockGitOperations) gitAdd(filename string) error {
	// Simulate git add
	return nil
}

func (m *MockGitOperations) gitCommit(message string) error {
	// Simulate git commit
	return nil
}

// Override the gitPush method to simulate a push without actually running the command
func (m *MockGitOperations) gitPush() error {
	fmt.Println("Simulating git push...")
	return nil // Simulate a successful push
}
