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
		{name: "subcommand help flag", args: []string{"scope", "--help"}, want: true},
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
