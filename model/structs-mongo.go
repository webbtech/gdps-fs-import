package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	ID          string             `bson:"_id"`
	AvgFuelCost float64            `bson:"avgFuelCost"`
	FuelSales   *FuelSales         `bson:"fuelSales"`
	ImportTS    int64              `bson:"importTS"`
	RecordDate  int                `bson:"recordDate"`
	StationID   primitive.ObjectID `bson:"stationID" json:"stationID"`
}

// FuelSalesImport struct
type FuelSalesImport struct {
	FuelCosts  *FuelCosts         `bson:"fuelCosts"`
	FuelSales  *FuelSales         `bson:"fuelSales"`
	FuelSums   *FuelSums          `bson:"fuelSums"`
	ImportTS   int64              `bson:"importTS"`
	RecordDate int                `bson:"recordDate"`
	StationID  primitive.ObjectID `bson:"stationID" json:"stationID"`
	Status     string             `bson:"status"`
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
	RecordDate  time.Time          `bson:"recordDate" json:"recordDate"`
	DispenserID primitive.ObjectID `bson:"dispenserID" json:"dispenserID"`
	Litres      float64            `bson:"litres" json:"litres"`
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
	ID    primitive.ObjectID   `bson:"_id"`
	Name  string               `bson:"name"`
	Nodes []primitive.ObjectID `bson:"nodes"`
}

// StationSales struct
type StationSales struct {
	RecordDate time.Time          `bson:"recordDate" json:"recordDate"`
	StationID  primitive.ObjectID `bson:"stationID" json:"stationID"`
	Fuel1      float64            `bson:"fuel1" json:"fuel1"`
	Fuel2      float64            `bson:"fuel2" json:"fuel2"`
	Fuel3      float64            `bson:"fuel3" json:"fuel3"`
	Fuel4      float64            `bson:"fuel4" json:"fuel4"`
	Fuel5      float64            `bson:"fuel5" json:"fuel5"`
	Fuel6      float64            `bson:"fuel6" json:"fuel6"`
	FuelCosts  *FuelCosts         `bson:"fuelCosts"`
}
