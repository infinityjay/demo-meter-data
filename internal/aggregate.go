package internal

import (
	"demo_meter_data/common"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type AggregateInfo struct {
	HouseholdID int
	Year        int
	Quarter     string
	Consumption int
}

func AggregateData(filename string, timeZone *time.Location) ([]AggregateInfo, int64, error) {
	start := time.Now().UnixMilli()
	var results []AggregateInfo
	aggregateMap := make(map[int]map[string]int)
	f, err := os.Open(filename)
	if err != nil {
		return nil, 0, fmt.Errorf("error opening %s: %v", filename, err)
	}
	lineNumber := 0
	defer f.Close()

	csvReader := csv.NewReader(f)
	// read the file line by line
	for {
		lineNumber++
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if lineNumber == 1 {
			continue
		}

		if len(row) != 3 {
			log.Printf("Invalid line length: %v, lineNumber: %v\n", row, lineNumber)
			continue
		}

		householdID, err := strconv.Atoi(row[0])
		if err != nil {
			log.Printf("Invalid householdID: %v, lineNumber: %v\n", err, lineNumber)
			continue
		}

		consumption, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			log.Printf("Invalid consumption: %v, lineNumber: %v, line: %v\n", err, lineNumber, row)
			continue
		}
		consumptionInt := int(math.Round(consumption * 100))
		// check if timestamp is valid
		year, quarter, err := common.ValidateTime(row[2], timeZone)
		if err != nil {
			log.Printf("Invalid timestamp: %v, lineNumber: %v\n", err, lineNumber)
			continue
		}
		quarterKey := fmt.Sprintf("%d-%s", year, quarter)
		// aggregate the consumption by quarter
		if _, ok := aggregateMap[householdID]; !ok {
			aggregateMap[householdID] = make(map[string]int)
		}
		aggregateMap[householdID][quarterKey] += consumptionInt
	}
	// get aggregated data from map
	for id, quarterData := range aggregateMap {
		var aggregatedData []AggregateInfo
		for quarterKey, consumption := range quarterData {
			year, quarter, err := parseQuarterKey(quarterKey)
			if err != nil {
				return nil, 0, fmt.Errorf("error parse quaterKey: %v", err)
			}
			aggregatedData = append(aggregatedData, AggregateInfo{
				HouseholdID: id,
				Year:        year,
				Quarter:     quarter,
				Consumption: consumption,
			})
		}
		results = append(results, aggregatedData...)
	}
	// sort results with id, year and quarter
	sort.Slice(results, func(i, j int) bool {
		if results[i].HouseholdID != results[j].HouseholdID {
			return results[i].HouseholdID < results[j].HouseholdID
		}
		if results[i].Year != results[j].Year {
			return results[i].Year < results[j].Year
		}
		return results[i].Quarter < results[j].Quarter
	})
	end := time.Now().UnixMilli()

	return results, end - start, nil
}

func parseQuarterKey(quarterKey string) (int, string, error) {
	parts := strings.Split(quarterKey, "-")
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid quarterKey format: %v", quarterKey)
	}
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", fmt.Errorf("invalid year in quarterKey: %v", err)
	}
	quarter := parts[1]
	if quarter != common.Q1 && quarter != common.Q2 && quarter != common.Q3 && quarter != common.Q4 {
		return 0, "", fmt.Errorf("invalid quarter in quarterKey: %v", err)
	}
	return year, quarter, nil
}
