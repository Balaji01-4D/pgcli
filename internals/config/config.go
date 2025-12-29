package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const filename = "config.toml"

type Config struct {
	prompt string `toml:"prompt"`
}

var DefaultConfig = Config{
	prompt: `\u@\h:\d> `,
}

func GetConfigDir() (string, error) {
	userdir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userdir, "pgxcli"), nil
}

func LoadConfig(path string) (config Config, err error) {
	var cfg Config
	_, err = toml.DecodeFile(path, &cfg)
	return cfg, err
}

func CheckConfigExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func SaveConfig(path string, cfg Config) error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}
	os.MkdirAll(dir, os.ModePerm)
	f, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return err
	}
	defer f.Close()

	enc := toml.NewEncoder(f)
	return enc.Encode(cfg)
}
