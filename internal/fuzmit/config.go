package fuzmit

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	EnvScope    = "FUZMIT_SCOPE"
	EnvGeoScope = "FUZMIT_JIRA_SCOPE"
	EnvNoEmojis = "FUZMIT_NO_EMOJIS"
)

// Config is the persisted default configuration.
type Config struct {
	Scope    bool `json:"scope"`
	GeoScope bool `json:"geoscope"`
	NoEmojis bool `json:"no_emojis"`
}

// Defaults is the resolved effective default set.
type Defaults struct {
	Scope    bool
	GeoScope bool
	NoEmojis bool
}

// ConfigPath returns the config file location.
func ConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}
	return filepath.Join(dir, "fuzmit", "config.json"), nil
}

// LoadConfig reads config from disk; missing file returns a zero config.
func LoadConfig() (Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return Config{}, err
	}

	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

// SaveConfig writes config to disk.
func SaveConfig(cfg Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	b = append(b, '\n')

	if err := os.WriteFile(path, b, 0o600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

// ResolveDefaults merges config defaults with environment overrides.
func ResolveDefaults(cfg Config, getenv func(string) string) Defaults {
	out := Defaults(cfg)

	if v, ok := parseBool(getenv(EnvScope)); ok {
		out.Scope = v
	}
	if v, ok := parseBool(getenv(EnvGeoScope)); ok {
		out.GeoScope = v
	}
	if v, ok := parseBool(getenv(EnvNoEmojis)); ok {
		out.NoEmojis = v
	}

	return out
}

// ParseToggleArg parses "on"/"off" style toggle args.
func ParseToggleArg(arg string) (bool, error) {
	if v, ok := parseBool(arg); ok {
		return v, nil
	}
	return false, fmt.Errorf("invalid value %q: expected on/off", arg)
}

func parseBool(raw string) (bool, bool) {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return false, false
	}

	switch raw {
	case "1", "true", "yes", "on", "y":
		return true, true
	case "0", "false", "no", "off", "n":
		return false, true
	}

	v, err := strconv.ParseBool(raw)
	if err != nil {
		return false, false
	}
	return v, true
}
