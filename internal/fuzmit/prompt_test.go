package fuzmit

import "testing"

func TestPickerTypesSortedByLengthThenAlpha(t *testing.T) {
	types := pickerTypes()
	if len(types) != len(SupportedTypes) {
		t.Fatalf("unexpected picker size: got %d want %d", len(types), len(SupportedTypes))
	}

	// pickerTypes is sorted by name length ascending, then alphabetically,
	// so shorter (more specific) names appear before longer ones.
	for i := 1; i < len(types); i++ {
		prev, cur := types[i-1], types[i]
		if len(prev.Name) > len(cur.Name) {
			t.Fatalf("picker order should be ascending by length: %q (len %d) > %q (len %d)",
				prev.Name, len(prev.Name), cur.Name, len(cur.Name))
		}
		if len(prev.Name) == len(cur.Name) && prev.Name > cur.Name {
			t.Fatalf("picker order should be alphabetical within same length: %q > %q",
				prev.Name, cur.Name)
		}
	}
}
