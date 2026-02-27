package fuzmit

import "testing"

func TestResolveDefaults_EnvOnly(t *testing.T) {
	env := map[string]string{
		EnvScope:    "true",
		EnvGeoScope: "false",
		EnvNoEmojis: "1",
	}

	getenv := func(k string) string { return env[k] }
	got := ResolveDefaults(getenv)

	if !got.Scope {
		t.Fatalf("scope should be true")
	}
	if got.GeoScope {
		t.Fatalf("geoscope should be false")
	}
	if !got.NoEmojis {
		t.Fatalf("no-emojis should be true")
	}
}

func TestResolveDefaults_DefaultsToFalse(t *testing.T) {
	got := ResolveDefaults(func(string) string { return "" })
	if got.Scope || got.GeoScope || got.NoEmojis {
		t.Fatalf("all defaults should be false: %#v", got)
	}
}

func TestResolveEnvSettings_InvalidValuesFallbackToDefault(t *testing.T) {
	env := map[string]string{
		EnvScope:    "maybe",
		EnvGeoScope: "TRUE",
	}
	settings := ResolveEnvSettings(func(k string) string { return env[k] })

	if len(settings) != 3 {
		t.Fatalf("expected 3 settings, got %d", len(settings))
	}

	if settings[0].Name != EnvScope || settings[0].Value || !settings[0].FromEnv || settings[0].ValidBool {
		t.Fatalf("unexpected scope setting: %#v", settings[0])
	}
	if settings[1].Name != EnvGeoScope || !settings[1].Value || !settings[1].FromEnv || !settings[1].ValidBool {
		t.Fatalf("unexpected geoscope setting: %#v", settings[1])
	}
	if settings[2].Name != EnvNoEmojis || settings[2].Value || settings[2].FromEnv || settings[2].ValidBool {
		t.Fatalf("unexpected no-emojis setting: %#v", settings[2])
	}
}
