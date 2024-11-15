package models

type RealEstate struct {
	SerialNumber    int
	ListYear        int
	DateRecorded    string
	Town            string
	Address         string
	AssessedValue   float64
	SaleAmount      float64
	SalesRatio      float64
	PropertyType    string
	ResidentialType string
	NonUseCode      string
	AssessorRemarks string
	OPMRemarks      string
	Location        string
}
