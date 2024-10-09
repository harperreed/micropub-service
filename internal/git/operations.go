// Package git provides functionality for managing git operations on blog posts.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"regexp"
	"gopkg.in/yaml.v2"
	"strings"
	"time"
	"log"
)

// RepoPath is the path to the content repository. It can be configured as needed.
var RepoPath = "./content"

// GitOperations interface defines the methods for git operations on blog posts.
type GitOperations interface {
	// CreatePost creates a new blog post with the given content.
	CreatePost(content map[string]interface{}) error

	// UpdatePost updates an existing blog post with new content.
	UpdatePost(content map[string]interface{}) error

	// DeletePost deletes an existing blog post.
	DeletePost(content map[string]interface{}) error

	// InitializeRepo initializes a new git repository for blog posts.
	InitializeRepo() error
}

// DefaultGitOperations is the default implementation of GitOperations.
type DefaultGitOperations struct{}

// GitOps is the global instance of GitOperations used throughout the package.
var GitOps GitOperations = &DefaultGitOperations{}

// CreatePost creates a new blog post with the given content.
// It takes a map containing the post properties and returns an error if the operation fails.
func (g *DefaultGitOperations) CreatePost(content map[string]interface{}) error {
	properties, ok := content["properties"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid properties")
	}

	var title, body string

	// Extract title
	if titleValue, ok := properties["title"]; ok {
		if titleSlice, ok := titleValue.([]interface{}); ok && len(titleSlice) > 0 {
			title, _ = titleSlice[0].(string)
		} else if titleStr, ok := titleValue.(string); ok {
			title = titleStr
		}
	}
	if title == "" {
		title = "Untitled Post"
	}

	// Extract content
	if contentValue, ok := properties["content"]; ok {
		if contentSlice, ok := contentValue.([]interface{}); ok && len(contentSlice) > 0 {
			body, _ = contentSlice[0].(string)
		} else if contentStr, ok := contentValue.(string); ok {
			body = contentStr
		}
	}
	if body == "" {
		return fmt.Errorf("missing content")
	}

	filename := fmt.Sprintf("%s-%s.md", time.Now().Format("2006-01-02"), sanitizeFilename(title))
	filePath := filepath.Join(RepoPath, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "---\ntitle: %s\ndate: %s\n---\n\n%s", title, time.Now().Format(time.RFC3339), body)
	if err != nil {
		return fmt.Errorf("failed to write content to file: %w", err)
	}

	if err := g.gitAdd(filename); err != nil {
		return fmt.Errorf("failed to git add: %w", err)
	}

	if err := g.gitCommit(fmt.Sprintf("Add post: %s", title)); err != nil {
		return fmt.Errorf("failed to git commit: %w", err)
	}

	if err := g.gitPush(); err != nil {
		return fmt.Errorf("failed to git push: %w", err)
	}

	// Set the URL in the content map
	content["url"] = fmt.Sprintf("/%s", filename)

	log.Printf("Created post: %s", filename)
	return nil
}

// UpdatePost updates an existing post with new content
func (g *DefaultGitOperations) UpdatePost(content map[string]interface{}) error {
	url, ok := content["url"].(string)
	if !ok {
		return fmt.Errorf("invalid URL")
	}

	properties, ok := content["properties"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid properties")
	}

	filename := filepath.Base(url)
	filePath := filepath.Join(RepoPath, filename)

	// Read existing content
	existingContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read existing file: %w", err)
	}

	frontmatter, oldContent, err := SplitFrontmatterAndContent(string(existingContent))
	if err != nil {
		return fmt.Errorf("failed to parse existing content: %w", err)
	}

	// Update title if provided
	if titleValue, ok := properties["title"]; ok {
		if titleArray, ok := titleValue.([]interface{}); ok && len(titleArray) > 0 {
			frontmatter["title"] = titleArray[0]
		} else if titleStr, ok := titleValue.(string); ok {
			frontmatter["title"] = titleStr
		}
	}

	// Update content if provided
	if contentValue, ok := properties["content"]; ok {
		if contentArray, ok := contentValue.([]interface{}); ok && len(contentArray) > 0 {
			oldContent = contentArray[0].(string)
		} else if contentStr, ok := contentValue.(string); ok {
			oldContent = contentStr
		}
	}

	updatedContent := CreateContentWithFrontmatter(frontmatter, oldContent)

	// Write updated content
	err = os.WriteFile(filePath, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated content: %w", err)
	}

	if err := g.gitAdd(filename); err != nil {
		return fmt.Errorf("failed to git add: %w", err)
	}

	if err := g.gitCommit(fmt.Sprintf("Update post: %s", filename)); err != nil {
		return fmt.Errorf("failed to git commit: %w", err)
	}

	if err := g.gitPush(); err != nil {
		return fmt.Errorf("failed to git push: %w", err)
	}

	log.Printf("Updated post: %s", url)
	return nil
}

// InitializeRepo initializes a new git repository for blog posts.
func (g *DefaultGitOperations) InitializeRepo() error {
	if _, err := os.Stat(RepoPath); os.IsNotExist(err) {
		err := os.MkdirAll(RepoPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create content directory: %w", err)
		}
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = RepoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	log.Printf("Initialized git repository at %s", RepoPath)
	return nil
}

// gitAdd adds a file to the git repository
func (g *DefaultGitOperations) gitAdd(filename string) error {
	cmd := exec.Command("git", "add", filename)
	cmd.Dir = RepoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git add: %w", err)
	}
	return nil
}

// gitCommit commits changes to the git repository
func (g *DefaultGitOperations) gitCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = RepoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git commit: %w", err)
	}
	return nil
}

// gitPush pushes changes to the remote git repository
func (g *DefaultGitOperations) gitPush() error {
	cmd := exec.Command("git", "push")
	cmd.Dir = RepoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git push: %w", err)
	}
	return nil
}

// Add this to your git/operations.go file
func SplitFrontmatterAndContent(content string) (map[string]interface{}, string, error) {
    parts := strings.SplitN(content, "---", 3)
    if len(parts) != 3 {
        return nil, "", fmt.Errorf("invalid frontmatter format")
    }

    var frontmatter map[string]interface{}
    err := yaml.Unmarshal([]byte(parts[1]), &frontmatter)
    if err != nil {
        return nil, "", err
    }

    return frontmatter, parts[2], nil
}

func CreateContentWithFrontmatter(frontmatter map[string]interface{}, content string) string {
    var sb strings.Builder
    sb.WriteString("---\n")

    // Sort the keys for consistent output
    var keys []string
    for k := range frontmatter {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    for _, k := range keys {
        v := frontmatter[k]
        sb.WriteString(fmt.Sprintf("%s: %v\n", k, v))
    }

    sb.WriteString("---\n")
    sb.WriteString(content)

    return sb.String()
}

func updateFrontMatter(content, key, value string) string {
    lines := strings.Split(content, "\n")
    inFrontMatter := false
    for i, line := range lines {
        if line == "---" {
            inFrontMatter = !inFrontMatter
            continue
        }
        if inFrontMatter && strings.HasPrefix(line, key+":") {
            lines[i] = fmt.Sprintf("%s: %s", key, value)
            break
        }
    }
    return strings.Join(lines, "\n")
}

func updateBody(content, newBody string) string {
    parts := strings.SplitN(content, "---", 3)
    if len(parts) < 3 {
        return content
    }
    return parts[0] + "---" + parts[1] + "---\n" + newBody
}

// This duplicate CreatePost method should be removed.
// The original CreatePost method is already defined earlier in the file.

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

    // Replace multiple hyphens with a single hyphen
    sanitized = regexp.MustCompile(`-+`).ReplaceAllString(sanitized, "-")

    // Trim any leading or trailing hyphens or underscores
    sanitized = strings.Trim(sanitized, "-_")

    return sanitized
}
