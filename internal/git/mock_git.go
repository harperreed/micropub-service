package git

import (
    "fmt"
)

type MockGitOperations struct {
    DefaultGitOperations
}

func (m *MockGitOperations) CreatePost(content map[string]interface{}) error {
    // Simulate creating a post
    return nil
}

func (m *MockGitOperations) UpdatePost(content map[string]interface{}) error {
    // Simulate updating a post
    return nil
}

func (m *MockGitOperations) DeletePost(content map[string]interface{}) error {
    // Simulate deleting a post
    return nil
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
