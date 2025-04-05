package internal

import (
	"testing"
)

func TestParseQuarterKey(t *testing.T) {
	tests := []struct {
		quarterKey      string
		expectedYear    int
		expectedQuarter string
		expectedResult  bool
	}{
		// Valid cases
		{"2024-Q1", 2024, "Q1", false},
		{"2024-Q2", 2024, "Q2", false},
		{"2024-Q3", 2024, "Q3", false},
		{"2024-Q4", 2024, "Q4", false},

		// Invalid quarter cases
		{"2023-Q5", 0, "", true},
		{"2023-Q", 0, "", true},

		// Invalid year cases
		{"invalid-Q1", 0, "", true},
		{"2023-QX", 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.quarterKey, func(t *testing.T) {
			year, quarter, err := parseQuarterKey(tt.quarterKey)

			// Compare the returned values with expected values
			if year != tt.expectedYear || quarter != tt.expectedQuarter {
				t.Errorf("expected %d-%s, got %d-%s", tt.expectedYear, tt.expectedQuarter, year, quarter)
			}

			// Check if an error is returned as expected
			if (err != nil) != tt.expectedResult {
				t.Errorf("got error: %v", err)
			}
		})
	}
}
