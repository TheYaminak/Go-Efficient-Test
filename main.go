package main

import (
	csvreader "GoEfficientTest/csvReader"
	"GoEfficientTest/dispatcher"
	"GoEfficientTest/handlers"
	"fmt"
)

func main() {

	// Start the servers
	handlers.StartServers()

	// Load records from CSV file
	records, err := csvreader.ReadCSV("Real_Estate_Sales_2001-2022_GL.csv")
	if err != nil {
		fmt.Printf("Error reading CSV file: %v\n", err)
		return
	}

	// Send requests to the servers
	dispatcher.SendRequestsToServers(records)
}
