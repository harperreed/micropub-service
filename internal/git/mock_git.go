package git

import (
    "fmt"
)

type MockGitOperations struct{}

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

func InitializeRepo() error {
    // Simulate initializing a repo
    return nil
}

func gitAdd(filename string) error {
    // Simulate git add
    return nil
}

func gitCommit(message string) error {
    // Simulate git commit
    return nil
}

func gitPush() error {
    // Simulate git push
    return nil
}
