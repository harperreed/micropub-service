package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	repoPath = "./content"
)

func initializeRepo() error {
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		err := os.MkdirAll(repoPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create content directory: %v", err)
		}
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %v", err)
	}

	return nil
}

func createPost(content map[string]interface{}) error {
	if err := initializeRepo(); err != nil {
		return err
	}

	title := content["title"].(string)
	body := content["content"].(string)

	filename := fmt.Sprintf("%s-%s.md", time.Now().Format("2006-01-02"), sanitizeFilename(title))
	filePath := filepath.Join(repoPath, filename)

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

	return nil
}

func gitAdd(filename string) error {
	cmd := exec.Command("git", "add", filename)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git add: %v", err)
	}
	return nil
}

func gitCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to git commit: %v", err)
	}
	return nil
}

func sanitizeFilename(filename string) string {
	// Implement a function to sanitize the filename
	// For simplicity, we'll just replace spaces with hyphens
	return filepath.Clean(strings.ReplaceAll(filename, " ", "-"))
}
