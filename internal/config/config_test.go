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
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "config_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Save the current working directory
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(oldWd)

	// Change to the temporary directory
	err = os.Chdir(tempDir)
	require.NoError(t, err)

	// Attempt to load the non-existent config file
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

func TestLoadConfigWithAdditionalFields(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "config_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a config file with additional fields
	configWithAdditionalFields := `{
		"gitRepoPath": "/path/to/repo",
		"additionalField1": "value1",
		"additionalField2": 42
	}`
	configPath := filepath.Join(tempDir, "config.json")
	err = os.WriteFile(configPath, []byte(configWithAdditionalFields), 0644)
	require.NoError(t, err)

	// Change to the temporary directory
	oldWd, _ := os.Getwd()
	err = os.Chdir(tempDir)
	require.NoError(t, err)
	defer os.Chdir(oldWd)

	// Load the config
	config, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "/path/to/repo", config.GitRepoPath)
}

func TestLoadConfigWithMissingRequiredFields(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "config_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a config file missing GitRepoPath
	configMissingRequired := `{
		"someOtherField": "value"
	}`
	configPath := filepath.Join(tempDir, "config.json")
	err = os.WriteFile(configPath, []byte(configMissingRequired), 0644)
	require.NoError(t, err)

	// Change to the temporary directory
	oldWd, _ := os.Getwd()
	err = os.Chdir(tempDir)
	require.NoError(t, err)
	defer os.Chdir(oldWd)

	// Attempt to load the config
	config, err := Load()
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "GitRepoPath is required")
}

func TestLoadConfigWithEnvironmentOverride(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "config_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a config file with a GitRepoPath
	configWithPath := `{
		"gitRepoPath": "/path/in/config"
	}`
	configPath := filepath.Join(tempDir, "config.json")
	err = os.WriteFile(configPath, []byte(configWithPath), 0644)
	require.NoError(t, err)

	// Set environment variable to override GitRepoPath
	envPath := "/path/from/env"
	os.Setenv("GIT_REPO_PATH", envPath)
	defer os.Unsetenv("GIT_REPO_PATH")

	// Change to the temporary directory
	oldWd, _ := os.Getwd()
	err = os.Chdir(tempDir)
	require.NoError(t, err)
	defer os.Chdir(oldWd)

	// Load the config
	config, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, envPath, config.GitRepoPath)
}
