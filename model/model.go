package model

import (
	"time"
)

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
	DateEnd    string `json:"dateEnd"`
	DateStart  string `json:"dateStart"`
}

// Request struct
type Request struct {
	DateEnd    time.Time
	DateStart  time.Time
	ExportType ExportType
}

// ErrorResponse struct
type ErrorResponse struct {
	Status  int    `json:"status"`
	Type    string `json:"type"`
	Message string `json:"message"`
}
