package main

import (
	csvreader "GoEfficientTest/csvReader"
	"GoEfficientTest/handlers"
	"GoEfficientTest/models"
	"GoEfficientTest/services"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
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
	SendRequestsToServers(records)
}

func SendRequestsToServers(records []models.RealEstate) {
	sequentialURL := "http://localhost:8080/process_sequential"
	concurrentURL := "http://localhost:8081/process_concurrent"

	var sequentialTotalTime, concurrentTotalTime time.Duration
	var sequentialRequestCount, concurrentRequestCount int

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine for sending requests to the sequential server
	go func() {
		defer wg.Done()
		for _, record := range records {
			// Validate data before sending
			if !services.ValidateData(record) {
				continue
			}

			// Marshal the record into JSON
			recordData, err := json.Marshal(record)
			if err != nil {
				fmt.Printf("Error marshaling record: %v\n", err)
				continue
			}

			// Send request to sequential server
			sequentialStart := time.Now()
			sequentialResp, err := http.Post(sequentialURL, "application/json", bytes.NewBuffer(recordData))
			if err != nil {
				fmt.Printf("Error sending request to sequential server: %v\n", err)
				continue
			}
			sequentialDuration := time.Since(sequentialStart)
			sequentialTotalTime += sequentialDuration
			sequentialRequestCount++
			sequentialResp.Body.Close()
		}
	}()

	// Goroutine for sending requests to the concurrent server
	go func() {
		defer wg.Done()
		for _, record := range records {
			// Validate data before sending
			if !services.ValidateData(record) {
				continue
			}

			// Marshal the record into JSON
			recordData, err := json.Marshal(record)
			if err != nil {
				fmt.Printf("Error marshaling record: %v\n", err)
				continue
			}

			// Send request to concurrent server
			concurrentStart := time.Now()
			concurrentResp, err := http.Post(concurrentURL, "application/json", bytes.NewBuffer(recordData))
			if err != nil {
				fmt.Printf("Error sending request to concurrent server: %v\n", err)
				continue
			}
			concurrentDuration := time.Since(concurrentStart)
			concurrentTotalTime += concurrentDuration
			concurrentRequestCount++
			concurrentResp.Body.Close()
		}
	}()

	// Wait for both goroutines to complete
	wg.Wait()

	// Calculate average response times
	var sequentialAverageTime, concurrentAverageTime time.Duration
	if sequentialRequestCount > 0 {
		sequentialAverageTime = sequentialTotalTime / time.Duration(sequentialRequestCount)
	}
	if concurrentRequestCount > 0 {
		concurrentAverageTime = concurrentTotalTime / time.Duration(concurrentRequestCount)
	}

	// Print metrics
	fmt.Printf("Total requests to sequential server: %d, Total time: %v, Average time: %v\n", sequentialRequestCount, sequentialTotalTime, sequentialAverageTime)
	fmt.Printf("Total requests to concurrent server: %d, Total time: %v, Average time: %v\n", concurrentRequestCount, concurrentTotalTime, concurrentAverageTime)
}
