package internal

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func GenerateCsv(data []AggregateInfo, outputFile string) error {
	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file, %v", err)
	}
	defer f.Close()

	csvwriter := csv.NewWriter(f)
	header := []string{"HouseholdID", "Year", "Quarter", "TotalConsumption"}
	csvwriter.Write(header)
	csvwriter.Flush()

	var lines [][]string
	for _, aggregate := range data {
		line := make([]string, 4)
		line[0] = strconv.Itoa(aggregate.HouseholdID)
		line[1] = strconv.Itoa(aggregate.Year)
		line[2] = aggregate.Quarter
		line[3] = fmt.Sprintf("%.2f", float64(aggregate.Consumption)/100.0)
		lines = append(lines, line)
	}
	err = csvwriter.WriteAll(lines)
	if err != nil {
		return fmt.Errorf("error writing to output file, %v", err)
	}
	return nil
}
