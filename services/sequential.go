package services

import (
	"GoEfficientTest/models"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

func ProcessRecordsWithService(records []models.RealEstate, serverURL string) time.Duration {
	totalStartTime := time.Now()
	var totalProcessingTime time.Duration

	for _, record := range records {
		// Send data to the server and wait for a response
		responseTime, err := sendRecordAndWait(record, serverURL)
		if err != nil {
			continue
		}
		totalProcessingTime += responseTime
	}

	totalDuration := time.Since(totalStartTime)
	return totalDuration
}

func sendRecordAndWait(record models.RealEstate, serverURL string) (time.Duration, error) {
	startTime := time.Now()

	data := []byte(fmt.Sprintf("SerialNumber: %d, Address: %s", record.SerialNumber, record.Address))
	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	processingDuration := time.Since(startTime)
	return processingDuration, nil
}
