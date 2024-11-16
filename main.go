package main

import (
	csvreader "GoEfficientTest/csvReader"
	"GoEfficientTest/dispatcher"
	"GoEfficientTest/handlers"
	"fmt"
	"time"
)

func main() {

	fmt.Println("Starting the servers...")
	// Start the servers
	go func() {
		handlers.StartServers()
	}()

	time.Sleep(5 * time.Second)
	// Load records from CSV file
	records, err := csvreader.ReadCSV("Real_Estate_Sales_2001-2022_GL.csv")
	if err != nil {
		fmt.Printf("Error reading CSV file: %v\n", err)
		return
	}

	// Send requests to the servers
	dispatcher.SendRequestsToServers(records)

	select {}
}
