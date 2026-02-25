package fuzmit

import "testing"

func TestFindCommitType(t *testing.T) {
	ct, ok := FindCommitType("fix")
	if !ok {
		t.Fatalf("expected fix to resolve")
	}
	if ct.Name != "fix" {
		t.Fatalf("expected fix, got %s", ct.Name)
	}

	if _, ok := FindCommitType("unknown"); ok {
		t.Fatalf("unknown type should not resolve")
	}
}
