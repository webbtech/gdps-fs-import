package model

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
