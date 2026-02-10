package excel

import "testing"

func TestColName(t *testing.T) {
	tests := []struct {
		index    int
		expected string
	}{
		{0, "A"},
		{1, "B"},
		{25, "Z"},
		{26, "AA"},
		{27, "AB"},
		{51, "AZ"},
		{52, "BA"},
	}
	for _, tt := range tests {
		got := colName(tt.index)
		if got != tt.expected {
			t.Errorf("colName(%d) = %s, want %s", tt.index, got, tt.expected)
		}
	}
}
