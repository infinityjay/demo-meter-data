package internal

import (
	"demo_meter_data/common"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"
)

func AggregateDataChunk(filename string, timeZone *time.Location) ([]AggregateInfo, int64, error) {
	start := time.Now().UnixMilli()
	f, err := os.Open(filename)
	if err != nil {
		return nil, 0, fmt.Errorf("error opening %s: %v", filename, err)
	}
	defer f.Close()
	// read file at once
	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, 0, fmt.Errorf("error reading CSV: %v", err)
	}
	if len(rows) <= 1 {
		return nil, 0, fmt.Errorf("no data rows found in file")
	}

	// remove the header
	dataRows := rows[1:]
	// Setup chunking
	numWorkers := runtime.NumCPU()
	fmt.Printf("CPU numbers: %v\n", numWorkers)
	chunkSize := (len(dataRows) + numWorkers - 1) / numWorkers
	chunks := splitIntoChunks(dataRows, chunkSize)

	var wg sync.WaitGroup
	var mu sync.Mutex
	resultsCh := make(chan map[int]map[string]int, numWorkers)

	for _, chunk := range chunks {
		wg.Add(1)
		go func(chunk [][]string) {
			defer wg.Done()
			aggregateMap := make(map[int]map[string]int)
			for _, row := range chunk {
				householdID, err := strconv.Atoi(row[0])
				if err != nil {
					log.Printf("Invalid householdID: %v, error: %v\n", row, err)
					continue
				}
				consumption, err := strconv.ParseFloat(row[1], 64)
				if err != nil {
					log.Printf("Invalid consumption: %v, error: %v\n", row, err)
					continue
				}
				consumptionInt := int(math.Round(consumption * 100))
				year, quarter, err := common.ValidateTime(row[2], timeZone)
				if err != nil {
					log.Printf("Invalid timestamp: %v, error: %v\n", row, err)
					continue
				}
				quarterKey := fmt.Sprintf("%d-%s", year, quarter)
				// add lock to avoid concurrent write to map
				mu.Lock()
				if _, ok := aggregateMap[householdID]; !ok {
					aggregateMap[householdID] = make(map[string]int)
				}
				aggregateMap[householdID][quarterKey] += consumptionInt
				mu.Unlock()
			}
			resultsCh <- aggregateMap
		}(chunk)
	}

	wg.Wait()
	close(resultsCh)

	// merge maps
	finalAgg := make(map[int]map[string]int)
	for part := range resultsCh {
		for id, quarterMap := range part {
			if _, ok := finalAgg[id]; !ok {
				finalAgg[id] = make(map[string]int)
			}
			for quarterKey, val := range quarterMap {
				finalAgg[id][quarterKey] += val
			}
		}
	}

	var results []AggregateInfo
	for id, quarterData := range finalAgg {
		for quarterKey, consumption := range quarterData {
			year, quarter, err := parseQuarterKey(quarterKey)
			if err != nil {
				return nil, 0, fmt.Errorf("error parse quarterKey: %v", err)
			}
			results = append(results, AggregateInfo{
				HouseholdID: id,
				Year:        year,
				Quarter:     quarter,
				Consumption: consumption,
			})
		}
	}

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

func splitIntoChunks(rows [][]string, chunkSize int) [][][]string {
	var chunks [][][]string
	for i := 0; i < len(rows); i += chunkSize {
		end := i + chunkSize
		if end > len(rows) {
			end = len(rows)
		}
		chunks = append(chunks, rows[i:end])
	}
	return chunks
}
