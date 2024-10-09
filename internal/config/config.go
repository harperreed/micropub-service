// Package config provides functionality for loading and managing application configuration.
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Config represents the application configuration.
type Config struct {
	// GitRepoPath is the path to the Git repository.
	GitRepoPath string `json:"gitRepoPath"`
}

// Load reads the configuration from a JSON file and returns a Config struct.
// It returns an error if the file cannot be read or parsed.
// Environment variable GIT_REPO_PATH can be used to override the GitRepoPath from the config file.
func Load() (*Config, error) {
	log.Println("Loading configuration from config.json")

	file, err := os.Open("config.json")
	if err != nil {
		log.Printf("Error opening config file: %v", err)
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("Error closing config file: %v", cerr)
		}
	}()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Printf("Error decoding config file: %v", err)
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// Check for environment variable override
	if envPath := os.Getenv("GIT_REPO_PATH"); envPath != "" {
		log.Println("Overriding GitRepoPath with environment variable")
		config.GitRepoPath = envPath
	}

	if config.GitRepoPath == "" {
		log.Println("GitRepoPath is empty in the config file and not set in environment")
		return nil, fmt.Errorf("GitRepoPath is required in the configuration or environment")
	}

	log.Println("Configuration loaded successfully")
	return &config, nil
}
