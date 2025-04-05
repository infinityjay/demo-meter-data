package main

import (
	"demo_meter_data/internal"
	"fmt"
	"time"
)

func main() {
	timeZone, _ := time.LoadLocation("Europe/Amsterdam")
	inputFile := "./data/household_consumption.csv"
	outputFile := "./data/quarterly_consumption.csv"
	// original process
	aggregatedData, processTime, err := internal.AggregateData(inputFile, timeZone)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("original processing time: %v\n", processTime)
	// chunk process
	aggregatedDataChunk, processTimeChunk, err := internal.AggregateDataChunk(inputFile, timeZone)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("original processing time: %v\n", processTimeChunk)
	// compare results
	if len(aggregatedData) != len(aggregatedDataChunk) {
		fmt.Println("the results of two function are different")
		return
	}
	for i, row := range aggregatedData {
		if row != aggregatedDataChunk[i] {
			fmt.Println("the results of two function are different")
			return
		}
	}

	err = internal.GenerateCsv(aggregatedData, outputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
}
