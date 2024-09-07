package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"gopkg.in/yaml.v2"
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

    properties, ok := content["properties"].(map[string]interface{})
    if !ok {
        return fmt.Errorf("invalid properties")
    }

    var title, body string

    if titleValue, ok := properties["title"]; ok {
        if titleArray, ok := titleValue.([]interface{}); ok && len(titleArray) > 0 {
            title, _ = titleArray[0].(string)
        } else if titleStr, ok := titleValue.(string); ok {
            title = titleStr
        }
    }

    if contentValue, ok := properties["content"]; ok {
        if contentArray, ok := contentValue.([]interface{}); ok && len(contentArray) > 0 {
            body, _ = contentArray[0].(string)
        } else if contentStr, ok := contentValue.(string); ok {
            body = contentStr
        }
    }

    if title == "" && body == "" {
        return fmt.Errorf("no updates provided")
    }

    filename := filepath.Base(url)
    filePath := filepath.Join(RepoPath, filename)

    // Read existing content
    existingContent, err := os.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("failed to read existing file: %v", err)
    }

    // Update content
    updatedContent := string(existingContent)
    if title != "" {
        updatedContent = updateFrontMatter(updatedContent, "title", title)
    }
    if body != "" {
        updatedContent = updateBody(updatedContent, body)
    }

    // Write updated content
    err = os.WriteFile(filePath, []byte(updatedContent), 0644)
    if err != nil {
        return fmt.Errorf("failed to write updated content: %v", err)
    }

    if err := gitAdd(filename); err != nil {
        return err
    }

    if err := gitCommit(fmt.Sprintf("Update post: %s", filename)); err != nil {
        return err
    }

    if err := gitPush(); err != nil {
        return err
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

    // Set the URL in the content map
    content["url"] = fmt.Sprintf("/%s", filename)

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
