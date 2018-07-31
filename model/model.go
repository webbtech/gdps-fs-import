package model

import "time"

// ExportType string
type ExportType string

// Export type constants
const (
	FuelType    ExportType = "fuel"
	PropaneType ExportType = "propane"
)

// RequestInput struct
type RequestInput struct {
	ExportType string `json:"exportType"`
	EndDate    string `json:"endDate"`
	StartDate  string `json:"startDate"`
}

// Request struct
type Request struct {
	EndDate    time.Time
	StartDate  time.Time
	ExportType ExportType
}
