package config

import (
	"os"
	"path"
	"pgcli/internals/logger"

	"github.com/pelletier/go-toml/v2"
)

type MainSection struct {
	MultiLine           bool     `toml:"multi_line"`
	DestructiveWarning  []string `toml:"destructive_warning"`
	LogFile             string   `toml:"log_file"`
	CasingFile          string   `toml:"casing_file"`
	GenerateCasingFile  bool     `toml:"generate_casing_file"`
	CaseColumnHeaders   bool     `toml:"case_column_headers"`
	HistoryFile         string   `toml:"history_file"`
	LogLevel            string   `toml:"log_level"`
	AsteriskColumnOrder string   `toml:"asterisk_column_order"`
	Timing              bool     `toml:"timing"`
	ShowBottomToolbar   bool     `toml:"show_bottom_toolbar"`
	TableFormat         string   `toml:"table_format"`
	RowLimit            int      `toml:"row_limit"`
	LessChatty          bool     `toml:"less_chatty"`
	Prompt              string   `toml:"prompt"`
	NullString          string   `toml:"null_string"`
	UseLocalTimezone    bool     `toml:"use_local_timezone"`
}

type DataFormats struct {
	Decimal string `toml:"decimal"`
	Float   string `toml:"float"`
}

type Config struct {
	Main              MainSection       `toml:"main"`
	NamedQueries      map[string]string `toml:"named queries"`
	AliasDSN          map[string]string `toml:"alias_dsn"`
	InitCommands      map[string]string `toml:"init-commands"`
	AliasDSNInitCmds  map[string]string `toml:"alias_dsn.init-commands"`
	DataFormats       DataFormats       `toml:"data_formats"`
	ColumnDateFormats map[string]string `toml:"column_date_formats"`
}

func Default() *Config {
	return &Config{
		Main: MainSection{
			MultiLine:           false,
			DestructiveWarning:  []string{"drop", "shutdown", "delete", "truncate", "alter", "update", "unconditional_update"},
			LogFile:             "default",
			CasingFile:          "default",
			GenerateCasingFile:  false,
			CaseColumnHeaders:   true,
			HistoryFile:         "default",
			LogLevel:            "INFO",
			AsteriskColumnOrder: "table_order",
			Timing:              true,
			ShowBottomToolbar:   true,
			TableFormat:         "psql",
			RowLimit:            1000,
			LessChatty:          false,
			Prompt:              "\\u@\\h:\\d> ",
			NullString:          "<null>",
			UseLocalTimezone:    true,
		},
		NamedQueries:      map[string]string{},
		AliasDSN:          map[string]string{},
		InitCommands:      map[string]string{},
		AliasDSNInitCmds:  map[string]string{},
		DataFormats:       DataFormats{Decimal: "", Float: ""},
		ColumnDateFormats: map[string]string{},
	}
}

func defaultTOML() (string, error) {
	cfg := Default()
	b, err := toml.Marshal(cfg)
	if err != nil {
		logger.Log.Error("Failed to marshal default config to TOML", "error", err)
		return "", err
	}

	out := string(b)
	return out, nil
}

func GetConfigDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(dir, "pgcli-go"), nil
}

func ensureConfigDirExists() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(configDir, 0o755)
}

func GetFilename(filename string) (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(configDir, filename), nil
}

func configFile() (string, error) { return GetFilename("config.toml") }

func WriteDefaultConfig() error {
	file, err := configFile()
	if err != nil {
		return err
	}

	eErr := ensureConfigDirExists()
	if eErr != nil {
		logger.Log.Error("Failed to ensure config directory exists", "error", eErr)
		return eErr
	}

	content, cErr := defaultTOML()
	if cErr != nil {
		logger.Log.Error("Failed to generate default TOML", "error", cErr)
		return cErr
	}

	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		logger.Log.Error("Failed to write default config to file", "error", err)
		return err
	}
	return nil
}

func Load() (*Config, error) {
	file, err := configFile()
	if err != nil {
		return nil, err
	}
	b, rErr := os.ReadFile(file)
	if rErr != nil {
		return nil, rErr
	}
	var cfg Config
	if uErr := toml.Unmarshal(b, &cfg); uErr != nil {
		return nil, uErr
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	file, err := configFile()
	if err != nil {
		return err
	}
	dir, _ := GetConfigDir()
	if mkErr := os.MkdirAll(dir, 0o755); mkErr != nil {
		return mkErr
	}
	b, mErr := toml.Marshal(c)
	if mErr != nil {
		return mErr
	}
	return os.WriteFile(file, b, 0o644)
}


