package services

import (
	"GoEfficientTest/models"
	"sync"
	"time"
)

func ProcessRecordsWithConcurrentService(records []models.RealEstate, serverURL string) time.Duration {
	totalStartTime := time.Now()
	var wg sync.WaitGroup
	workerCount := 10
	jobs := make(chan models.RealEstate, len(records))

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for record := range jobs {
				_, err := sendRecordAndWait(record, serverURL)
				if err != nil {
					continue
				}
			}
		}()
	}

	for _, record := range records {
		jobs <- record
	}
	close(jobs)

	wg.Wait()
	totalDuration := time.Since(totalStartTime)
	return totalDuration
}
