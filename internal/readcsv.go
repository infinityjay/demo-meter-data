package internal

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type MeterInfo struct {
	HouseholdID int
	Consumption float32
	Timestamp   int64
}

func ParseData(filename string) (*MeterInfo, error) {
	results := new(MeterInfo)
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %v", filename, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
	}
	return results, nil
}
