package export

import "fmt"

func (e *Exporter) propane() (err error) {

	fmt.Println("fueltype:", e.Request.ExportType)
	return nil
}
