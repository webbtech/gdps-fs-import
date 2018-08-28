package export

import (
	"github.com/pulpfree/gales-fuelsale-export/config"
	"github.com/pulpfree/gales-fuelsale-export/model"
)

const timeForm = "2006-01-02"

// Exporter struct
type Exporter struct {
	Request *model.Request
	cfg     *config.Config
}

// New function
func New(r *model.Request, cfg *config.Config) *Exporter {
	e := &Exporter{Request: r, cfg: cfg}
	return e
}

// Process request function
func (e *Exporter) Process() (res *model.DnImportRes, err error) {

	switch e.Request.ExportType {
	case model.FuelType:
		res, err = e.fuel()
	case model.PropaneType:
		res, err = e.propane()
	}

	return res, err
}
