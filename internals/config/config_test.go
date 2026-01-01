package config_test

import (
	"os"
	path "path/filepath"
	"pgcli/internals/config"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestLoadConfig_ValidConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := path.Join(tempDir, "config.toml")
	configContent := `prompt = "\\u@\\h:\\d> "`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	cfg, err := config.LoadConfig(configPath)
	assert.NoError(t, err)
	assert.Equal(t, "\\u@\\h:\\d> ", cfg.Prompt)
}

func TestLoadConfig_InvalidConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := path.Join(tempDir, "config.toml")
	invalidContent := `prompt = "\\u@\\h:\\d> ` // Missing closing quote

	err := os.WriteFile(configPath, []byte(invalidContent), 0644)
	assert.NoError(t, err)

	_, err = config.LoadConfig(configPath)
	assert.Error(t, err)
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := config.LoadConfig("non_existent_config.toml")
	assert.Error(t, err)
}

func TestGetConfigDir_ShouldReturnValidPath(t *testing.T) {
	tempdir := t.TempDir()
	os.Setenv("HOME", tempdir)
	configDir, err := config.GetConfigDir()
	assert.NoError(t, err)
	expectedPath := path.Join(tempdir, ".config", "pgxcli")
	assert.Equal(t, expectedPath, configDir)
} 

