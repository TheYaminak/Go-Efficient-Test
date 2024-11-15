package handlers

import (
	"GoEfficientTest/models"
	"GoEfficientTest/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func SequentialHandler(w http.ResponseWriter, r *http.Request) {
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

	startTime := time.Now()
	for _, record := range records {
		if !services.ValidateData(record) {
			continue
		}
		services.AdjustValues(&record)
	}
	processingDuration := time.Since(startTime)

	response := fmt.Sprintf("Processed %d records sequentially in %v", len(records), processingDuration)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
