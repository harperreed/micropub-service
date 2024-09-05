package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	repoPath  = "./content"
	remoteURL = "https://github.com/your-username/your-repo.git"
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

	if err := gitPush(); err != nil {
		return err
	}

	return nil
}

func updatePost(content map[string]interface{}) error {
	if err := initializeRepo(); err != nil {
		return err
	}

	url := content["url"].(string)
	title := content["title"].(string)
	body := content["content"].(string)

	filename := filepath.Base(url)
	filePath := filepath.Join(repoPath, filename)

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

func deletePost(content map[string]interface{}) error {
	if err := initializeRepo(); err != nil {
		return err
	}

	url := content["url"].(string)
	filename := filepath.Base(url)
	filePath := filepath.Join(repoPath, filename)

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

func gitPush() error {
	branchName := fmt.Sprintf("post-%d", time.Now().Unix())

	// Create a new branch
	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create new branch: %v", err)
	}

	// Push the new branch to the remote repository
	cmd = exec.Command("git", "push", "-u", "origin", branchName)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push new branch: %v", err)
	}

	// Create a pull request (this is a placeholder, as creating a PR typically requires using the GitHub API)
	fmt.Printf("New branch '%s' has been pushed. Please create a pull request manually.\n", branchName)

	return nil
}

func sanitizeFilename(filename string) string {
	// Implement a function to sanitize the filename
	// For simplicity, we'll just replace spaces with hyphens
	return filepath.Clean(strings.ReplaceAll(filename, " ", "-"))
}
