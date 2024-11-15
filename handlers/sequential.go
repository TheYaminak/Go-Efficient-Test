package handlers

import (
	"GoEfficientTest/models"
	"GoEfficientTest/services"
	"fmt"
	"net/http"
	"time"
)

func SequentialHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	record := getRecordFromRequest(r)
	if record == nil {
		http.Error(w, "Failed to parse record", http.StatusBadRequest)
		return
	}

	// Start processing data sequentially
	processingStart := time.Now()
	processingDuration := services.ProcessRecordsSequential([]models.RealEstate{*record})
	services.AdjustValues(record)
	averageSaleAmount, averageAssessedValue, residentialTypeCount, propertyTypeCount := services.CalculateStatistics([]models.RealEstate{*record})
	salesRatio := services.CalculateSalesRatios([]models.RealEstate{*record})
	//townStats := services.AnalyzeTownData([]models.RealEstate{*record})
	processingTime := time.Since(processingStart)

	// Update global metrics
	mutex.Lock()
	sequentialTotalTime += processingTime
	sequentialRequestCount++
	sequentialProcessingTime += processingDuration
	sequentialSaleStats = append(sequentialSaleStats, averageSaleAmount)
	sequentialSalesRatios = append(sequentialSalesRatios, salesRatio)
	//sequentialTownData = append(sequentialTownData, townStats)
	mutex.Unlock()

	// Prepare response
	response := fmt.Sprintf(
		"Processed record sequentially\nProcessing Time: %v\nAverage Sale Amount: %.2f\nAverage Assessed Value: %.2f\nSales Ratio: %.2f\nResidential Type Count: %v\nProperty Type Count: %v\n",
		processingTime, averageSaleAmount, averageAssessedValue, salesRatio, residentialTypeCount, propertyTypeCount,
	)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
