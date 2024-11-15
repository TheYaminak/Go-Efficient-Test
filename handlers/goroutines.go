package handlers

import (
	"GoEfficientTest/models"
	"GoEfficientTest/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func ConcurrentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var records []models.RealEstate
	err = json.Unmarshal(body, &records)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	// Send initial response that the data has been validated
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data validated successfully, processing started."))

	// Start processing data concurrently
	go func() {
		var wg sync.WaitGroup
		startTime := time.Now()
		for _, record := range records {
			if !services.ValidateData(record) {
				continue
			}
			wg.Add(1)
			go func(rec models.RealEstate) {
				defer wg.Done()
				services.AdjustValues(&rec)
			}(record)
		}
		wg.Wait()
		processingDuration := time.Since(startTime)
		fmt.Printf("Processed %d records concurrently in %v\n", len(records), processingDuration)
	}()
}
