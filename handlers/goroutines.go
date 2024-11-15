package handlers

import (
	"GoEfficientTest/models"
	"GoEfficientTest/services"
	"fmt"
	"net/http"
	"time"
)

func ConcurrentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	record := getRecordFromRequest(r)
	if record == nil {
		http.Error(w, "Failed to parse record", http.StatusBadRequest)
		return
	}

	// Start processing data concurrently
	processingStart := time.Now()
	processingDuration := services.ProcessRecordsConcurrent([]models.RealEstate{*record})
	services.AdjustValues(record)
	averageSaleAmount, averageAssessedValue, residentialTypeCount, propertyTypeCount := services.CalculateStatistics([]models.RealEstate{*record})
	salesRatio := services.CalculateSalesRatios([]models.RealEstate{*record})
	//townStats := services.AnalyzeTownData([]models.RealEstate{*record})
	processingTime := time.Since(processingStart)

	// Update global metrics
	mutex.Lock()
	concurrentTotalTime += processingTime
	concurrentRequestCount++
	concurrentProcessingTime += processingDuration
	concurrentSaleStats = append(concurrentSaleStats, averageSaleAmount)
	concurrentSalesRatios = append(concurrentSalesRatios, salesRatio)
	//concurrentTownData = append(concurrentTownData, townStats)
	mutex.Unlock()

	// Prepare response
	response := fmt.Sprintf("Processed record concurrently\nProcessing Time: %v\nAverage Sale Amount: %.2f\nAverage Assessed Value: %.2f\nSales Ratio: %.2f\nResidential Type Count: %v\nProperty Type Count: %v\n", processingTime, averageSaleAmount, averageAssessedValue, salesRatio, residentialTypeCount, propertyTypeCount)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
