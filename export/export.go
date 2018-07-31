package export

import (
	"errors"

	"github.com/pulpfree/gales-fuelsale-export/config"
	"github.com/pulpfree/gales-fuelsale-export/model"
	log "github.com/sirupsen/logrus"
)

// Exporter struct
type Exporter struct {
	Request *model.Request
	cfg     *config.Config
}

// New function
func New(r *model.Request, cfg *config.Config) *Exporter {

	e := &Exporter{Request: r, cfg: cfg}

	// fmt.Printf("cfg: %+v\n", cfg)

	return e
}

// Process request function
func (e *Exporter) Process() (err error) {

	// fmt.Printf("request: %+v\n", e.Request)
	log.Error("Failed to fetch fuel sales count")
	// log.Errorf("Failed to fetch fuel sales count: %s", err)

	switch e.Request.ExportType {
	case model.FuelType:
		err = e.fuel()
	case model.PropaneType:
		err = e.propane()
	default:
		err = errors.New("Invalid fuel type requested")
	}
	if err != nil {
		return err
	}

	return nil
}
