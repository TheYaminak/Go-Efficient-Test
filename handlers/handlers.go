package handlers

import (
	"GoEfficientTest/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	sequentialTotalTime      time.Duration
	concurrentTotalTime      time.Duration
	sequentialRequestCount   int
	concurrentRequestCount   int
	sequentialProcessingTime time.Duration
	concurrentProcessingTime time.Duration
	sequentialSaleStats      []float64
	concurrentSaleStats      []float64
	sequentialSalesRatios    []float64
	concurrentSalesRatios    []float64
	// sequentialTownData       map[string]struct {
	// 	totalSaleAmount    float64
	// 	totalAssessedValue float64
	// 	numProperties      int
	// }
	// concurrentTownData map[string]struct {
	// 	totalSaleAmount    float64
	// 	totalAssessedValue float64
	// 	numProperties      int
	// }
	mutex sync.Mutex
)

func StartServers() {
	http.HandleFunc("/process_sequential", SequentialHandler)
	go func() {
		http.ListenAndServe(":8081", nil)
	}()

	http.HandleFunc("/process_concurrent", ConcurrentHandler)
	http.ListenAndServe(":8082", nil)
}

func getRecordFromRequest(r *http.Request) *models.RealEstate {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	defer r.Body.Close()

	var record models.RealEstate
	err = json.Unmarshal(body, &record)
	if err != nil {
		return nil
	}

	return &record
}

func PrintMetrics() {
	mutex.Lock()
	defer mutex.Unlock()

	var sequentialAverageTime, concurrentAverageTime, sequentialAverageProcessingTime, concurrentAverageProcessingTime time.Duration
	if sequentialRequestCount > 0 {
		sequentialAverageTime = sequentialTotalTime / time.Duration(sequentialRequestCount)
		sequentialAverageProcessingTime = sequentialProcessingTime / time.Duration(sequentialRequestCount)
	}
	if concurrentRequestCount > 0 {
		concurrentAverageTime = concurrentTotalTime / time.Duration(concurrentRequestCount)
		concurrentAverageProcessingTime = concurrentProcessingTime / time.Duration(concurrentRequestCount)
	}

	fmt.Printf("Total requests to sequential server: %d, Total time: %v, Average time: %v, Average processing time: %v\n", sequentialRequestCount, sequentialTotalTime, sequentialAverageTime, sequentialAverageProcessingTime)
	fmt.Printf("Total requests to concurrent server: %d, Total time: %v, Average time: %v, Average processing time: %v\n", concurrentRequestCount, concurrentTotalTime, concurrentAverageTime, concurrentAverageProcessingTime)

	fmt.Printf("Sequential Sale Statistics: %v\n", sequentialSaleStats)
	fmt.Printf("Concurrent Sale Statistics: %v\n", concurrentSaleStats)
	fmt.Printf("Sequential Sales Ratios: %v\n", sequentialSalesRatios)
	fmt.Printf("Concurrent Sales Ratios: %v\n", concurrentSalesRatios)
	// fmt.Printf("Sequential Town Data: %v\n", sequentialTownData)
	// fmt.Printf("Concurrent Town Data: %v\n", concurrentTownData)
}
