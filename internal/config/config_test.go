package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "config_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a valid config file
	validConfig := Config{GitRepoPath: "/path/to/repo"}
	configPath := filepath.Join(tempDir, "config.json")
	configData, err := json.Marshal(validConfig)
	require.NoError(t, err)
	err = os.WriteFile(configPath, configData, 0644)
	require.NoError(t, err)

	// Test successful config loading
	t.Run("SuccessfulLoad", func(t *testing.T) {
		oldWd, _ := os.Getwd()
		err = os.Chdir(tempDir)
		require.NoError(t, err)
		defer os.Chdir(oldWd)

		config, err := Load()
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, validConfig.GitRepoPath, config.GitRepoPath)
	})

	// Test loading with empty GitRepoPath
	t.Run("EmptyGitRepoPath", func(t *testing.T) {
		emptyConfig := Config{GitRepoPath: ""}
		emptyConfigData, err := json.Marshal(emptyConfig)
		require.NoError(t, err)
		err = os.WriteFile(configPath, emptyConfigData, 0644)
		require.NoError(t, err)

		oldWd, _ := os.Getwd()
		err = os.Chdir(tempDir)
		require.NoError(t, err)
		defer os.Chdir(oldWd)

		config, err := Load()
		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "GitRepoPath is required")
	})
}

func TestLoadNonExistentFile(t *testing.T) {
	// Set the current working directory to a non-existent directory
	err := os.Chdir("/non/existent/directory")
	require.NoError(t, err)
	defer os.Chdir("/")

	config, err := Load()
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "failed to open config file")
}

func TestLoadInvalidJSON(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "config_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create an invalid JSON config file
	configPath := filepath.Join(tempDir, "config.json")
	err = os.WriteFile(configPath, []byte("invalid json"), 0644)
	require.NoError(t, err)

	// Set the current working directory to the temp directory
	oldWd, _ := os.Getwd()
	err = os.Chdir(tempDir)
	require.NoError(t, err)
	defer os.Chdir(oldWd)

	config, err := Load()
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "failed to decode config file")
}

func TestConfigStruct(t *testing.T) {
	config := Config{GitRepoPath: "/test/repo/path"}
	assert.Equal(t, "/test/repo/path", config.GitRepoPath)

	// Test JSON marshaling and unmarshaling
	jsonData, err := json.Marshal(config)
	assert.NoError(t, err)

	var unmarshaledConfig Config
	err = json.Unmarshal(jsonData, &unmarshaledConfig)
	assert.NoError(t, err)
	assert.Equal(t, config, unmarshaledConfig)
}
