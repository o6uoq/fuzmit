package main

import "testing"

func TestHelpRequested(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{name: "root help command", args: []string{"help"}, want: true},
		{name: "root long help flag", args: []string{"--help"}, want: true},
		{name: "root short help flag", args: []string{"-h"}, want: true},
		{name: "subcommand help flag", args: []string{"env", "--help"}, want: true},
		{name: "regular command", args: []string{"--type", "fix", "-m", "trim panic"}, want: false},
		{name: "description value help", args: []string{"-m", "help"}, want: false},
		{name: "separator before flag", args: []string{"--", "--help"}, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := helpRequested(tc.args); got != tc.want {
				t.Fatalf("helpRequested(%v)=%v want %v", tc.args, got, tc.want)
			}
		})
	}
}

func TestNormalizeArgs(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "scope value as separate token",
			in:   []string{"--type", "fix", "--scope", "auth", "-m", "trim panic"},
			want: []string{"--type", "fix", "--scope=auth", "-m", "trim panic"},
		},
		{
			name: "short scope value as separate token",
			in:   []string{"-t", "fix", "-s", "auth", "-m", "trim panic"},
			want: []string{"-t", "fix", "-s=auth", "-m", "trim panic"},
		},
		{
			name: "scope prompt remains bare",
			in:   []string{"--type", "fix", "--scope", "-m", "trim panic"},
			want: []string{"--type", "fix", "--scope", "-m", "trim panic"},
		},
		{
			name: "do not normalize after separator",
			in:   []string{"--", "--scope", "auth"},
			want: []string{"--", "--scope", "auth"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeArgs(tc.in)
			if len(got) != len(tc.want) {
				t.Fatalf("normalizeArgs(%v) len=%d want %d (%v)", tc.in, len(got), len(tc.want), got)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("normalizeArgs(%v)[%d]=%q want %q (full=%v)", tc.in, i, got[i], tc.want[i], got)
				}
			}
		})
	}
}
