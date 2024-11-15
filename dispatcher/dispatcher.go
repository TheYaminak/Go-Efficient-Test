package dispatcher

import (
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

func SendRequestsToServers(records []models.RealEstate) {
	sequentialURL := "http://localhost:8081/process_sequential"
	concurrentURL := "http://localhost:8082/process_concurrent"

	var sequentialTotalTime, concurrentTotalTime time.Duration
	var sequentialRequestCount, concurrentRequestCount int

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		totalRecords := len(records)

		for index, record := range records {
			if !services.ValidateData(record) {
				continue
			}

			recordData, err := json.Marshal(record)
			if err != nil {
				fmt.Printf("Error marshaling record: %v\n", err)
				continue
			}

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

			// Calculate progress percentage and print
			progress := float64(index+1) / float64(totalRecords) * 100
			fmt.Printf("\rSequential Progress: %.2f%%", progress)
		}
		fmt.Println() // Print newline after completion
	}()

	go func() {
		defer wg.Done()
		totalRecords := len(records)

		for index, record := range records {
			if !services.ValidateData(record) {
				continue
			}

			recordData, err := json.Marshal(record)
			if err != nil {
				fmt.Printf("Error marshaling record: %v\n", err)
				continue
			}

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

			// Calculate progress percentage and print
			progress := float64(index+1) / float64(totalRecords) * 100
			fmt.Printf("\rConcurrent Progress: %.2f%%", progress)
		}
		fmt.Println() // Print newline after completion
	}()

	wg.Wait()

	var sequentialAverageTime, concurrentAverageTime time.Duration
	if sequentialRequestCount > 0 {
		sequentialAverageTime = sequentialTotalTime / time.Duration(sequentialRequestCount)
	}
	if concurrentRequestCount > 0 {
		concurrentAverageTime = concurrentTotalTime / time.Duration(concurrentRequestCount)
	}

	fmt.Printf("Total requests to sequential server: %d, Total time: %v, Average time: %v\n", sequentialRequestCount, sequentialTotalTime, sequentialAverageTime)
	fmt.Printf("Total requests to concurrent server: %d, Total time: %v, Average time: %v\n", concurrentRequestCount, concurrentTotalTime, concurrentAverageTime)

	handlers.PrintMetrics()
}
