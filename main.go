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
	aggregatedData, err := internal.AggregateData(inputFile, timeZone)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = internal.GenerateCsv(aggregatedData, outputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
}
