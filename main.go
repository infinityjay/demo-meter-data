package main

import (
	"demo_meter_data/internal"
	"fmt"
	"github.com/spf13/pflag"
	"time"
)

func main() {
	var inputFile = pflag.StringP("inputFile", "i", "./data/household_consumption.csv", "input file path")
	var outputFile = pflag.StringP("outputFile", "o", "./data/quarterly_consumption.csv", "output file path")
	var chunk = pflag.BoolP("chunk", "c", false, "use chunk function or not")
	pflag.Parse()

	timeZone, _ := time.LoadLocation("Europe/Amsterdam")

	if *chunk {
		// chunk process
		aggregatedDataChunk, processTimeChunk, err := internal.AggregateDataChunk(*inputFile, timeZone)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("original processing time: %v\n", processTimeChunk)

		// generate csv file
		err = internal.GenerateCsv(aggregatedDataChunk, *outputFile)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		// original process
		aggregatedData, processTime, err := internal.AggregateData(*inputFile, timeZone)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("original processing time: %v\n", processTime)

		// generate csv file
		err = internal.GenerateCsv(aggregatedData, *outputFile)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
