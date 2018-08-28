package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
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

// ================ gales-sales DB structs ================ //

// FuelCosts struct
type FuelCosts struct {
	Fuel1 float64 `bson:"fuel_1" json:"fuel1"`
	Fuel2 float64 `bson:"fuel_2" json:"fuel2"`
	Fuel3 float64 `bson:"fuel_3" json:"fuel3"`
	Fuel4 float64 `bson:"fuel_4" json:"fuel4"`
	Fuel5 float64 `bson:"fuel_5" json:"fuel5"`
	Fuel6 float64 `bson:"fuel_6" json:"fuel6"`
}

// FuelSales struct
type FuelSales struct {
	NL   float64 `bson:"NL" json:"NL"`
	SNL  float64 `bson:"SNL" json:"SNL"`
	DSL  float64 `bson:"DSL" json:"DSL"`
	CDSL float64 `bson:"CDSL" json:"CDSL"`
	PROP float64 `bson:"PROP" json:"PROP"`
}

// FuelSalesExport struct
type FuelSalesExport struct {
	ID          string        `bson:"_id"`
	AvgFuelCost float64       `bson:"avgFuelCost"`
	FuelSales   *FuelSales    `bson:"fuelSales"`
	ImportTS    int64         `bson:"importTS"`
	RecordDate  int           `bson:"recordDate"`
	StationID   bson.ObjectId `bson:"stationID" json:"stationID"`
}

// FuelSalesImport struct
type FuelSalesImport struct {
	FuelCosts  *FuelCosts    `bson:"fuelCosts"`
	FuelSales  *FuelSales    `bson:"fuelSales"`
	FuelSums   *FuelSums     `bson:"fuelSums"`
	ImportTS   int64         `bson:"importTS"`
	RecordDate int           `bson:"recordDate"`
	StationID  bson.ObjectId `bson:"stationID" json:"stationID"`
	Status     string        `bson:"status"`
}

// FuelSums struct
type FuelSums struct {
	Fuel1 float64 `bson:"fuel1" json:"fuel1"`
	Fuel2 float64 `bson:"fuel2" json:"fuel2"`
	Fuel3 float64 `bson:"fuel3" json:"fuel3"`
	Fuel4 float64 `bson:"fuel4" json:"fuel4"`
	Fuel5 float64 `bson:"fuel5" json:"fuel5"`
	Fuel6 float64 `bson:"fuel6" json:"fuel6"`
}

// ImportLog struct
type ImportLog struct {
	DateFrom   int        `bson:"dateFrom"`
	DateTo     int        `bson:"dateTo"`
	ImportTS   int64      `bson:"importTS"`
	ImportType ExportType `bson:"importType"`
}

// PropaneSale struct
type PropaneSale struct {
	RecordDate  time.Time     `bson:"recordDate" json:"recordDate"`
	DispenserID bson.ObjectId `bson:"dispenserID" json:"dispenserID"`
	Litres      float64       `bson:"litres" json:"litres"`
}

// PropaneSaleExport struct
type PropaneSaleExport struct {
	ImportTS   int64   `bson:"importTS"`
	Litres     float64 `bson:"litres" json:"litres"`
	RecordDate int     `bson:"recordDate"`
	TankID     int     `bson:"tankID" json:"TankID"`
}

// StationNodes struct
type StationNodes struct {
	ID    bson.ObjectId   `bson:"_id"`
	Name  string          `bson:"name"`
	Nodes []bson.ObjectId `bson:"nodes"`
}

// StationSales struct
type StationSales struct {
	RecordDate time.Time     `bson:"recordDate" json:"recordDate"`
	StationID  bson.ObjectId `bson:"stationID" json:"stationID"`
	Fuel1      float64       `bson:"fuel1" json:"fuel1"`
	Fuel2      float64       `bson:"fuel2" json:"fuel2"`
	Fuel3      float64       `bson:"fuel3" json:"fuel3"`
	Fuel4      float64       `bson:"fuel4" json:"fuel4"`
	Fuel5      float64       `bson:"fuel5" json:"fuel5"`
	Fuel6      float64       `bson:"fuel6" json:"fuel6"`
	FuelCosts  *FuelCosts    `bson:"fuelCosts"`
}

// ================ DynamoDB structs ================ //

// DnStation struct
type DnStation struct {
	ID         string `json:"ID"`
	Name       string `json:"Name"`
	RefStation string `json:"RefStation"`
}

// DnFuelSales struct
type DnFuelSales struct {
	AvgFuelCost float64    `json:"AvgFuelCost"`
	Date        int        `json:"Date"`
	ImportTS    int64      `json:"ImportTS"`
	Sales       *FuelSales `json:"Sales"`
	StationID   string     `json:"StationID"`
	YearWeek    int        `json:"YearWeek"`
}

// DnFuelPrice struct
type DnFuelPrice struct {
	Date      int     `json:"Date"`
	Price     float64 `json:"Price"`
	StationID string  `json:"StationID"`
	YearWeek  int     `json:"YearWeek"`
}

// DnImportRes struct
type DnImportRes struct {
	DateEnd        string `json:"DateEnd"`
	DateStart      string `json:"DateStart"`
	ImportDate     string `json:"ImportDate"`
	ImportTS       int64  `json:"ImportTS"`
	ImportType     string `json:"ImportType"`
	RecordQuantity int    `json:"RecordQty"`
}

// DnPropaneSales struct
type DnPropaneSales struct {
	Date     int     `json:"Date"`
	ImportTS int64   `json:"ImportTS"`
	Sales    float64 `json:"Sales"`
	TankID   int     `json:"TankID"`
	Year     int     `json:"Year"`
	YearWeek int     `json:"YearWeek"`
}
