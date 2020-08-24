package export

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/pulpfree/gsales-fs-export/model"
	"github.com/pulpfree/gsales-fs-export/model/dynamo"
	"github.com/pulpfree/gsales-fs-export/model/mongo"
)

func (e *Exporter) fuel() (res *model.DnImportRes, err error) {

	// Set MongoDB connection
	mongo, err := mongo.NewDB(e.cfg.GetMongoConnectURL(), e.cfg.MongoDBName)
	if err != nil {
		log.Errorf("Error connecting to mongo: %s", err)
		return res, err
	}
	defer mongo.Close()

	// Set DynamoDB connection
	dynamo, err := dynamo.NewDB(e.cfg.Dynamo)
	if err != nil {
		log.Errorf("Error connecting to dynamo: %s", err)
		return res, err
	}

	t := time.Now()
	res = &model.DnImportRes{
		DateEnd:    e.Request.DateEnd.Format(timeForm),
		DateStart:  e.Request.DateStart.Format(timeForm),
		ImportDate: t.Format(timeForm),
		ImportTS:   t.Unix(),
		ImportType: string(e.Request.ExportType),
	}

	// Create and fetch mongo fuel sales records
	err = mongo.CreateFuelSales(e.Request)
	if err != nil {
		log.Errorf("Error creating fuel sales: %s", err)
		return res, err
	}
	sales, err := mongo.FetchExportedFuelSales(e.Request)
	if err != nil {
		log.Errorf("Error fetching fuel sales: %s", err)
		return res, err
	}
	if len(sales) <= 0 {
		err = errors.New("Error fetching exported fuel sales")
		log.Error(err)
		return res, err
	}

	res.RecordQuantity = len(sales)

	err = dynamo.CreateFuelSalesRecords(sales, res)
	if err != nil {
		log.Errorf("Error creating dynamo sales records: %s", err)
		return res, err
	}

	return res, err
}
