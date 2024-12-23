package services

import (
	"GoEfficientTest/models"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

func ValidateData(re models.RealEstate) bool {
	if re.SaleAmount < 0 || re.AssessedValue < 0 {
		return false
	}
	return true
}

func AdjustValues(re *models.RealEstate) {
	factor := rand.Float64()*(1.2-0.8) + 0.8 // Factor between 0.8 and 1.2
	re.AssessedValue *= factor

	yearDifference := time.Now().Year() - re.ListYear
	appreciationRate := 0.03 // Assume an annual appreciation rate of 3%
	if yearDifference > 0 {
		re.SaleAmount *= (1 + appreciationRate*float64(yearDifference))
	}
}

func CalculateStatistics(records []models.RealEstate) (float64, float64, map[string]int, map[string]int) {
	var totalSaleAmount float64
	var totalAssessedValue float64
	residentialTypeCount := make(map[string]int)
	propertyTypeCount := make(map[string]int)

	for _, re := range records {
		totalSaleAmount += re.SaleAmount
		totalAssessedValue += re.AssessedValue
		residentialTypeCount[re.ResidentialType]++
		propertyTypeCount[re.PropertyType]++
	}

	numRecords := float64(len(records))
	averageSaleAmount := totalSaleAmount / numRecords
	averageAssessedValue := totalAssessedValue / numRecords

	return averageSaleAmount, averageAssessedValue, residentialTypeCount, propertyTypeCount
}

func CalculateSalesRatios(records []models.RealEstate) float64 {
	var totalSalesRatio float64
	var count int

	for _, re := range records {
		if re.SalesRatio > 0 {
			totalSalesRatio += re.SalesRatio
			count++
		}
	}

	if count > 0 {
		return totalSalesRatio / float64(count)
	}
	return 0.0
}

//Not used for now
// func AnalyzeTownData(records []models.RealEstate) map[string]struct {
// 	totalSaleAmount    float64
// 	totalAssessedValue float64
// 	numProperties      int
// } {

// 	townStats := make(map[string]struct {
// 		totalSaleAmount    float64
// 		totalAssessedValue float64
// 		numProperties      int
// 	})

// 	for _, re := range records {
// 		stats := townStats[re.Town]
// 		stats.totalSaleAmount += re.SaleAmount
// 		stats.totalAssessedValue += re.AssessedValue
// 		stats.numProperties++
// 		townStats[re.Town] = stats
// 	}

// 	return townStats
// }

func ProcessRecordsSequential(records []models.RealEstate) time.Duration {
	startTime := time.Now()
	for _, re := range records {
		processRecord(re)
	}
	return time.Since(startTime)
}

func ProcessRecordsConcurrent(records []models.RealEstate) time.Duration {
	startTime := time.Now()
	var wg sync.WaitGroup

	for _, re := range records {
		wg.Add(1)
		go func(record models.RealEstate) {
			defer wg.Done()
			processRecord(record)
		}(re)
	}

	wg.Wait()
	return time.Since(startTime)
}

func processRecord(re models.RealEstate) time.Duration {
	startTime := time.Now()

	if !ValidateData(re) {
		return time.Since(startTime)
	}

	AdjustValues(&re)

	return time.Since(startTime)
}

func ExportToCSV(data []float64, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"Value"}); err != nil {
		return fmt.Errorf("failed to write headers: %v", err)
	}

	for _, value := range data {
		record := []string{fmt.Sprintf("%.2f", value)}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %v", err)
		}
	}

	return nil
}
