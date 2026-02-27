package fuzmit

import (
	"fmt"
	"strings"
)

const (
	EnvScope    = "FUZMIT_SCOPE"
	EnvGeoScope = "FUZMIT_JIRA_SCOPE"
	EnvNoEmojis = "FUZMIT_NO_EMOJIS"
)

// Defaults is the resolved effective default set.
type Defaults struct {
	Scope    bool
	GeoScope bool
	NoEmojis bool
}

// EnvSetting is the resolved value for a supported FUZMIT_* variable.
type EnvSetting struct {
	Name      string
	Value     bool
	Raw       string
	FromEnv   bool
	ValidBool bool
}

// ResolveDefaults reads defaults from environment variables only.
func ResolveDefaults(getenv func(string) string) Defaults {
	settings := ResolveEnvSettings(getenv)
	out := Defaults{}
	for _, s := range settings {
		switch s.Name {
		case EnvScope:
			out.Scope = s.Value
		case EnvGeoScope:
			out.GeoScope = s.Value
		case EnvNoEmojis:
			out.NoEmojis = s.Value
		}
	}
	return out
}

// ResolveEnvSettings returns supported environment settings and their resolved values.
func ResolveEnvSettings(getenv func(string) string) []EnvSetting {
	keys := []string{EnvScope, EnvGeoScope, EnvNoEmojis}
	settings := make([]EnvSetting, 0, len(keys))
	for _, key := range keys {
		raw := strings.TrimSpace(getenv(key))
		value, ok := parseBool(raw)
		settings = append(settings, EnvSetting{
			Name:      key,
			Value:     value,
			Raw:       raw,
			FromEnv:   raw != "",
			ValidBool: ok,
		})
	}
	return settings
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

	return false, false
}

func describeEnvSetting(setting EnvSetting) string {
	if !setting.FromEnv {
		return fmt.Sprintf("%t (default)", setting.Value)
	}
	if !setting.ValidBool {
		return fmt.Sprintf("%t (default; invalid env value %q)", setting.Value, setting.Raw)
	}
	return fmt.Sprintf("%t (from env: %q)", setting.Value, setting.Raw)
}
