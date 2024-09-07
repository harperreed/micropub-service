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