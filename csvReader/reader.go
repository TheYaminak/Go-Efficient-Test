package csvreader

import (
	"GoEfficientTest/models"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func ReadCSV(filename string) ([]models.RealEstate, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("could not read header: %w", err)
	}

	var records []models.RealEstate
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("could not read record: %w", err)
		}

		serialNumber, _ := strconv.Atoi(record[0])
		listYear, _ := strconv.Atoi(record[1])
		assessedValue, _ := strconv.ParseFloat(record[5], 64)
		saleAmount, _ := strconv.ParseFloat(record[6], 64)
		salesRatio, _ := strconv.ParseFloat(record[7], 64)

		re := models.RealEstate{
			SerialNumber:    serialNumber,
			ListYear:        listYear,
			DateRecorded:    record[2],
			Town:            record[3],
			Address:         record[4],
			AssessedValue:   assessedValue,
			SaleAmount:      saleAmount,
			SalesRatio:      salesRatio,
			PropertyType:    record[8],
			ResidentialType: record[9],
			NonUseCode:      record[10],
			AssessorRemarks: record[11],
			OPMRemarks:      record[12],
			Location:        record[13],
		}
		records = append(records, re)
	}

	return records, nil
}
