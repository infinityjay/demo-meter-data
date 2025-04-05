package common

import (
	"testing"
	"time"
)

func TestValidateTime(t *testing.T) {
	// Define a valid timezone (e.g., UTC)
	loc, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		t.Fatalf("Failed to load location: %v", err)
	}

	tests := []struct {
		name            string
		timeString      string
		expectedYear    int
		expectedQuarter string
		expectError     bool
	}{
		{
			name:            "Valid timestamp - January 15, 2021",
			timeString:      "1610694000",
			expectedYear:    2021,
			expectedQuarter: "Q1",
			expectError:     false,
		},
		{
			name:            "Valid timestamp - November 15, 2025",
			timeString:      "1732047600",
			expectedYear:    2024,
			expectedQuarter: "Q4",
			expectError:     false,
		},
		{
			name:            "Invalid timestamp format",
			timeString:      "invalid_timestamp",
			expectedYear:    0,
			expectedQuarter: "",
			expectError:     true,
		},
		{
			name:            "Timestamp out of range - December 31 1999",
			timeString:      "946656000",
			expectedYear:    0,
			expectedQuarter: "",
			expectError:     true,
		},
		{
			name:            "Timestamp out of range - January 1 2100",
			timeString:      "4102444800",
			expectedYear:    0,
			expectedQuarter: "",
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			year, quarter, err := ValidateTime(tt.timeString, loc)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
			if !tt.expectError {
				if year != tt.expectedYear {
					t.Errorf("expected year: %v, got: %v", tt.expectedYear, year)
				}
				if quarter != tt.expectedQuarter {
					t.Errorf("expected quarter: %v, got: %v", tt.expectedQuarter, quarter)
				}
			}
		})
	}
}
