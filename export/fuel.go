package export

import (
	"fmt"
	"log"

	"github.com/pulpfree/gales-fuelsale-export/model/mongo"
)

func (e *Exporter) fuel() (err error) {

	// fmt.Println("fueltype:", e.Request.ExportType)
	// fmt.Println("startDate:", e.Request.StartDate)
	// fmt.Println("cfg:", e.cfg)
	// fmt.Printf("cfg: %+v\n", e.cfg)

	db, err := mongo.NewDB(e.cfg.GetMongoConnectURL())
	if err != nil {
		log.Fatalf("Error connecting to mongo: %s", err)
	}

	fmt.Printf("db: %+v\n", db)

	return nil
}
