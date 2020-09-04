package validators

import (
	"errors"
	"time"

	"github.com/pulpfree/gsales-fs-export/model"
)

// Time form constant
const (
	timeShortForm  = "20060102"
	timeRecordForm = "2006-01-02"
)

// Date function
func Date(dateInput string) (time.Time, error) {

	date, err := time.Parse(timeRecordForm, dateInput)
	if err != nil {
		return date, err
	}

	// Ensure date is not future dated
	today := time.Now()
	if today.Unix() < date.Unix() {
		return date, errors.New("Invalid date. Date cannot be future, must be less than current date")
	}

	return date, err
}

// Fuel function
func Fuel(exportInput string) (model.ExportType, error) {
	switch exportInput {
	case "fuel":
		return model.FuelType, nil
	case "propane":
		return model.PropaneType, nil
	default:
		return "", errors.New("Invalid export type provided")
	}
}

// RequestVars function
func RequestVars(r *model.RequestInput) (res *model.Request, err error) {

	res = new(model.Request)
	res.ExportType, err = Fuel(r.ExportType)
	if err != nil {
		return res, err
	}

	res.DateStart, err = Date(r.DateStart)
	if err != nil {
		return res, err
	}
	res.DateEnd, err = Date(r.DateEnd)
	if err != nil {
		return res, err
	}

	return res, nil
}
