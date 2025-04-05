package internal

import (
	"os"
	"testing"
	"time"
)

func TestAggregateData(t *testing.T) {
	tests := []struct {
		name           string
		inputCSV       string
		expectedResult []AggregateInfo
		expectedError  error
	}{
		{
			name: "Valid data aggregation",
			inputCSV: `householdID,consumption,timestamp
1,10.33,1675596305
1,10.6,1675682705
2,20.66,1676460305
2,2.4,1675682705
1,10.5,1684149905
invalid,100,1684149905
1,invalid,1684149905
2,2.34,1683285905`,
			expectedResult: []AggregateInfo{
				{HouseholdID: 1, Year: 2023, Quarter: "Q1", Consumption: 2093},
				{HouseholdID: 1, Year: 2023, Quarter: "Q2", Consumption: 1050},
				{HouseholdID: 2, Year: 2023, Quarter: "Q1", Consumption: 2306},
				{HouseholdID: 2, Year: 2023, Quarter: "Q2", Consumption: 234},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a temporary file with the input CSV data
			tmpFile, err := os.CreateTemp("", "*.csv")
			if err != nil {
				t.Fatalf("Error creating temporary file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			// write input CSV to the temporary file
			_, err = tmpFile.WriteString(tt.inputCSV)
			if err != nil {
				t.Fatalf("Error writing to temporary file: %v", err)
			}
			tmpFile.Close()

			timeZone, _ := time.LoadLocation("Europe/Amsterdam")
			result, _, err := AggregateData(tmpFile.Name(), timeZone)
			if (err != nil) != (tt.expectedError != nil) {
				t.Errorf("expected error: %v, got error: %v", tt.expectedError, err)
			}
			// check result
			if len(result) != len(tt.expectedResult) {
				t.Errorf("expected result: %v, got result: %v", tt.expectedResult, result)
			} else {
				for i, r := range result {
					if r != tt.expectedResult[i] {
						t.Errorf("expected result[%d]: %v, got result[%d]: %v", i, tt.expectedResult[i], i, r)
					}
				}
			}
		})
	}
}

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
