package validators

import (
	"testing"
	"time"

	"github.com/pulpfree/gales-fuelsale-export/model"
	"github.com/stretchr/testify/assert"
)

// TestValidDate test for valid date format only
func TestValidDate(t *testing.T) {
	date, err := Date("2018-05-01")

	assert.NoError(t, err)
	assert.IsType(t, time.Time{}, date)
}

// TestInValidDate test for invalid date format
func TestInValidDate(t *testing.T) {
	_, err := Date("2018-051")

	assert.Error(t, err)
}

// TestInvalidFutureDate test for invalid future date
func TestInvalidFutureDate(t *testing.T) {

	today := time.Now()
	futureDate := today.Add(time.Hour * 24 * 2).Format(timeRecordForm)

	_, err := Date(futureDate)

	assert.Error(t, err)
}

// TestValidExportType function
func TestValidExportType(t *testing.T) {

	var et model.ExportType

	tp, err := Fuel("fuelSales")
	assert.NoError(t, err)
	assert.IsType(t, et, tp)
}

// TestInValidExportType function
func TestInValidExportType(t *testing.T) {

	_, err := Fuel("fulType")

	assert.Error(t, err)
}

// TestValidRequest function
func TestValidRequest(t *testing.T) {

	testVars := &model.RequestInput{
		StartDate:  "2018-01-01",
		EndDate:    "2018-02-28",
		ExportType: "fuelSales",
	}

	res, err := RequestVars(testVars)

	var reqTp *model.Request

	assert.NoError(t, err)
	assert.IsType(t, reqTp, res)
}
