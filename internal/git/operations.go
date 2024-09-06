package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var RepoPath = "./content" // You might want to make this configurable

// GitOperations interface defines the methods for git operations
type GitOperations interface {
	CreatePost(content map[string]interface{}) error
	UpdatePost(content map[string]interface{}) error
	DeletePost(content map[string]interface{}) error
}


// DefaultGitOperations is the default implementation of GitOperations
type DefaultGitOperations struct{}

var GitOps GitOperations = &DefaultGitOperations{}

func (g *DefaultGitOperations) UpdatePost(content map[string]interface{}) error {
	url, ok := content["url"].(string)
	if !ok {
		return fmt.Errorf("invalid URL")
	}

	title, ok := content["title"].(string)
	if !ok {
		return fmt.Errorf("invalid title")
	}

	body, ok := content["content"].(string)
	if !ok {
		return fmt.Errorf("invalid content")
	}

	filename := filepath.Base(url)
	filePath := filepath.Join(RepoPath, filename)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "---\ntitle: %s\ndate: %s\n---\n\n%s", title, time.Now().Format(time.RFC3339), body)
	if err != nil {
		return fmt.Errorf("failed to write content to file: %v", err)
	}

	if err := gitAdd(filename); err != nil {
		return err
	}

	if err := gitCommit(fmt.Sprintf("Update post: %s", title)); err != nil {
		return err
	}

	if err := gitPush(); err != nil {
		return err
	}

	return nil
}

func (g *DefaultGitOperations) CreatePost(content map[string]interface{}) error {
	title := content["title"].(string)
	body := content["content"].(string)

	filename := fmt.Sprintf("%s-%s.md", time.Now().Format("2006-01-02"), sanitizeFilename(title))
	filePath := filepath.Join(RepoPath, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "---\ntitle: %s\ndate: %s\n---\n\n%s", title, time.Now().Format(time.RFC3339), body)
	if err != nil {
		return fmt.Errorf("failed to write content to file: %v", err)
	}

	if err := gitAdd(filename); err != nil {
		return err
	}

	if err := gitCommit(fmt.Sprintf("Add post: %s", title)); err != nil {
		return err
	}

	if err := gitPush(); err != nil {
		return err
	}

	return nil
}

func InitializeRepo() error {
	if _, err := os.Stat(RepoPath); os.IsNotExist(err) {
		err := os.MkdirAll(RepoPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create content directory: %v", err)
		}
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = RepoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %v", err)
	}

	return nil
}

func gitAdd(filename string) error {
	cmd := exec.Command("git", "add", filename)
	cmd.Dir = RepoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git add: %v", err)
	}
	return nil
}

func gitCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = RepoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git commit: %v", err)
	}
	return nil
}

func gitPush() error {
	cmd := exec.Command("git", "push")
	cmd.Dir = RepoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git push: %v", err)
	}
	return nil
}



func (g *DefaultGitOperations) DeletePost(content map[string]interface{}) error {
	url := content["url"].(string)
	filename := filepath.Base(url)
	filePath := filepath.Join(RepoPath, filename)

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	if err := gitAdd(filename); err != nil {
		return err
	}

	if err := gitCommit(fmt.Sprintf("Delete post: %s", filename)); err != nil {
		return err
	}

	if err := gitPush(); err != nil {
		return err
	}

	return nil
}

func sanitizeFilename(filename string) string {
	// Replace spaces with hyphens
	sanitized := strings.ReplaceAll(filename, " ", "-")

	// Remove any characters that aren't alphanumeric, hyphen, or underscore
	sanitized = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return -1
	}, sanitized)

	// Convert to lowercase
	sanitized = strings.ToLower(sanitized)

	// Trim any leading or trailing hyphens or underscores
	sanitized = strings.Trim(sanitized, "-_")

	return sanitized
}
