package internal

import (
	"testing"
	"time"
)

func TestAggregateChunk(t *testing.T) {
	// Load the timezone
	timeZone, _ := time.LoadLocation("Europe/Amsterdam")
	inputFile := "../data/household_consumption.csv"

	// original function
	aggregatedData, _, err := AggregateData(inputFile, timeZone)
	if err != nil {
		t.Fatalf("Error in AggregateData: %v", err)
	}

	// chunk function
	aggregatedDataChunk, _, err := AggregateDataChunk(inputFile, timeZone)
	if err != nil {
		t.Fatalf("Error in AggregateDataChunk: %v", err)
	}

	// compare results
	if len(aggregatedData) != len(aggregatedDataChunk) {
		t.Errorf("The results of two functions are different: different lengths")
		return
	}
	for i, row := range aggregatedData {
		if row != aggregatedDataChunk[i] {
			t.Errorf("The results of two functions are different at index %d", i)
			return
		}
	}
}
