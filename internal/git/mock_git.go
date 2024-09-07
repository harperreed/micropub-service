package git

import (
    "fmt"
    "time"
    "path/filepath"
    "os"
)

type MockGitOperations struct {
    DefaultGitOperations
}

// Add these methods to your MockGitOperations struct

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
