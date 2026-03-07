package fuzmit

import "testing"

func TestPickerTypesDisplaysAlphabeticalAtEmptyQuery(t *testing.T) {
	types := pickerTypes()
	if len(types) != len(SupportedTypes) {
		t.Fatalf("unexpected picker size: got %d want %d", len(types), len(SupportedTypes))
	}

	// pickerTypes is sorted alphabetically so the picker always starts in
	// conventional type-name order.
	for i := 1; i < len(types); i++ {
		if types[i-1].Name > types[i].Name {
			t.Fatalf("picker order should be ascending alphabetical: %q > %q", types[i-1].Name, types[i].Name)
		}
	}
}
