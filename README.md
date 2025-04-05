# demo-meter-data

## Introduction
This is a simple implementation to aggregate the Meter Data into fiscal quarters.
The input and output files are located in the `./data` folder.

## Execute
You can execute the program by following command and then you can check the results in the `./data` folder.
```bash
# default input file is ./data/household_consumption.csv
go run main.go
```

I also implement the optimized version for the process logic. You can generate the test file and try following command.
```bash
# generate test file
go run ./data/expand.go

# specify input or output file
go run main.go -i ./data/expanded.csv -o ./data/quarterly_expanded.csv

# specify use chunk function, default not use
go run main.go -c true 

# check parameters
go run main.go --help
```

## Test
I provide unit tests for the processing functions, and can run the tests by following commands
```bash
go test ./common -v 

go test ./internal -v
```

## Optimization

### Scenario

There could be an extremely large dataset where processing each row of consumption data is time-consuming.

### Analysis

For this example, the bottleneck is in the processing part. The results  will only rely on the unique household number, and  
the process of writing the result to the csv file does not need to be optimized. And the main goal of the optimization 
is to improve the time efficiency of line by line data processing, and we assume the memory usage is not limited.

### Solutions

1. Parallel processing

We can set a workers pool according to the CPU cores and create multiple goroutines to process the data. We can 
read the csv file line by line and then assign the line to different goroutines to process the logic. And collect
the data from different goroutines.

2. Multi-chunks processing

We can separate the file into different chunks at the beginning, and then process each chunk with individual goroutines.
Similar to the solution 1, we can collect the data from different chunk results.

3. Use column-based databases

If the output of the csv file is not an immediate request, we can also store the data into column-based databases like DuckDB. 
I once implemented a performance test program to compare the performance of specific query between row-based database (MySql) 
and column-based database(DuckDB).

The DuckDB will highly outperform on the tasks like aggregating data based on time, you can find my test from [here](https://github.com/infinityjay/ADM2024/blob/main/ADM2024-A2-s4162315/code.ipynb).

### Comparison

I implement the solution 2 Multi-chunks processing. I set the number of workers to the number of CPUs and split the 
rows evenly into multiple chunks. I generate a 1.57GB file by duplicate the original file for 500 times through the 
file `./data/expand.go`.

Then I run the main function, it turns out that the original function processing time is 37590ms, and the optimized 
function processing time is 20543 which is almost a 45% speedup. Use the following command separately.

```bash
# generate test file
go run ./data/expand.go

# original function
go run main.go -i ./data/expanded.csv -o ./data/quarterly_expanded.csv

# use chunked function
go run main.go -i ./data/expanded.csv -o ./data/quarterly_expanded.csv -c true 
```

But the result is still not so ideal, and the read file part is also time-consuming, we can consider using streaming
 file read.



