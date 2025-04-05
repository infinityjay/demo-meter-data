package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	inputPath := "./data/household_consumption.csv"
	outputPath := "./data/expanded.csv"
	multiplier := 500

	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic(fmt.Errorf("failed to open input file: %v", err))
	}
	defer inputFile.Close()

	var lines []string
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("error reading lines: %v", err))
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic(fmt.Errorf("failed to create output file: %v", err))
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	for i := 0; i < multiplier; i++ {
		for _, line := range lines {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				panic(fmt.Errorf("error writing line: %v", err))
			}
		}
	}
	writer.Flush()

	fmt.Printf("Expanded file written to %s (%d× original)\n", outputPath, multiplier)
}
