package fuzmit

import "testing"

func TestResolveDefaults_EnvOverridesConfig(t *testing.T) {
	cfg := Config{Scope: true, GeoScope: false, NoEmojis: false}
	env := map[string]string{
		EnvScope:    "false",
		EnvGeoScope: "true",
		EnvNoEmojis: "1",
	}

	getenv := func(k string) string { return env[k] }
	got := ResolveDefaults(cfg, getenv)

	if got.Scope {
		t.Fatalf("scope should be false")
	}
	if !got.GeoScope {
		t.Fatalf("geoscope should be true")
	}
	if !got.NoEmojis {
		t.Fatalf("no-emojis should be true")
	}
}

func TestParseToggleArg(t *testing.T) {
	v, err := ParseToggleArg("on")
	if err != nil || !v {
		t.Fatalf("on should parse true: v=%t err=%v", v, err)
	}

	v, err = ParseToggleArg("off")
	if err != nil || v {
		t.Fatalf("off should parse false: v=%t err=%v", v, err)
	}

	if _, err := ParseToggleArg("maybe"); err == nil {
		t.Fatalf("maybe should fail")
	}
}
