package dispatcher

import (
	"GoEfficientTest/handlers"
	"GoEfficientTest/models"
	"GoEfficientTest/services"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func SendRequestsToServers(records []models.RealEstate) {

	client := &http.Client{}
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

			progress := float64(index+1) / float64(totalRecords) * 100
			fmt.Printf("\rSequential Progress: %.2f%%", progress)
		}
		fmt.Println()
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

			request, err := http.NewRequest("POST", concurrentURL, bytes.NewBuffer(recordData))
			if err != nil {
				fmt.Printf("Error creating request: %v\n", err)
				continue
			}
			request.Header.Set("Content-Type", "application/json")

			concurrentStart := time.Now()
			concurrentResp, err := client.Do(request)
			if err != nil {
				fmt.Printf("Error sending request to concurrent server: %v\n", err)
				continue
			}
			concurrentDuration := time.Since(concurrentStart)
			concurrentTotalTime += concurrentDuration
			concurrentRequestCount++

			// Close the response body properly to free up resources
			_, _ = io.Copy(io.Discard, concurrentResp.Body) // Asegurarse de leer todo el cuerpo de la respuesta
			concurrentResp.Body.Close()

			progress := float64(index+1) / float64(totalRecords) * 100
			fmt.Printf("\rConcurrent Progress: %.2f%%", progress)
		}
		fmt.Println()
	}()
	wg.Wait()

	handlers.PrintMetrics()
}
