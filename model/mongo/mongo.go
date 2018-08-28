package mongo

import (
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pulpfree/gales-fuelsale-export/config"
	"github.com/pulpfree/gales-fuelsale-export/model"
)

// DB struct
type DB struct {
	session *mgo.Session
}

// DB Constants
const (
	DBSales         = "gales-sales"
	colFuelSales    = "fuel-sales"
	colFSImport     = "fuel-sales-import"
	colFSExport     = "fuel-sales-export"
	colImportLog    = "import-log"
	colPSExport     = "propane-sales-export"
	colSales        = "sales"
	colStationNodes = "station-nodes"
)

// Time format constants
const (
	timeShortForm = "20060102"
)

// ==================== Exported methods ==================== //

// NewDB connection function
func NewDB(connection string) (*DB, error) {

	s, err := mgo.Dial(connection)
	if err != nil {
		return nil, err
	}

	return &DB{
		session: s,
	}, err
}

// CreateFuelSales function
func (db *DB) CreateFuelSales(req *model.Request) (err error) {

	sales, err := db.fetchFuelSales(req)
	if err != nil {
		return err
	}

	ts, err := db.persistFuelSales(sales)
	if err != nil {
		return err
	}

	err = db.createImportLog(req, ts)
	if err != nil {
		return err
	}

	err = db.compileFuelSales()
	if err != nil {
		return err
	}

	err = db.removeImportedFuelSales()
	if err != nil {
		return err
	}

	return err
}

// CreatePropaneSales function
func (db *DB) CreatePropaneSales(req *model.Request) (err error) {

	sales, err := db.fetchPropaneSales(req)
	if err != nil {
		return err
	}

	ts, err := db.persistPropaneSales(sales)
	if err != nil {
		return err
	}

	err = db.createImportLog(req, ts)
	if err != nil {
		return err
	}

	return err
}

// FetchExportedFuelSales method
func (db *DB) FetchExportedFuelSales(req *model.Request) (res []*model.FuelSalesExport, err error) {

	s := db.getFreshSession()
	defer s.Close()

	// fetch previously exported records by date range
	col := s.DB(DBSales).C(colFSExport)
	stDte, _ := strconv.Atoi(req.DateStart.Format(timeShortForm))
	enDte, _ := strconv.Atoi(req.DateEnd.Format(timeShortForm))
	col.Find(bson.M{"recordDate": bson.M{"$gte": stDte, "$lte": enDte}}).All(&res)

	return res, err
}

// FetchExportedPropaneSales method
func (db *DB) FetchExportedPropaneSales(req *model.Request) (res []*model.PropaneSaleExport, err error) {

	s := db.getFreshSession()
	defer s.Close()

	// fetch previously exported records by date range
	col := s.DB(DBSales).C(colPSExport)
	stDte, _ := strconv.Atoi(req.DateStart.Format(timeShortForm))
	enDte, _ := strconv.Atoi(req.DateEnd.Format(timeShortForm))
	col.Find(bson.M{"recordDate": bson.M{"$gte": stDte, "$lte": enDte}}).All(&res)

	return res, err
}

// ==================== FuelSales methods ==================== //

func (db *DB) fetchFuelSales(req *model.Request) (ss []*model.StationSales, err error) {

	s := db.getFreshSession()
	defer s.Close()

	col := s.DB(DBSales).C(colSales)
	match := bson.M{
		"$match": bson.M{"recordDate": bson.M{"$gte": req.DateStart, "$lte": req.DateEnd}},
	}

	group := bson.M{
		"$group": bson.M{
			"_id":       bson.M{"recordDate": "$recordDate", "stationID": "$stationID"},
			"fuel1":     bson.M{"$sum": "$salesSummary.fuel.fuel_1.litre"},
			"fuel2":     bson.M{"$sum": "$salesSummary.fuel.fuel_2.litre"},
			"fuel3":     bson.M{"$sum": "$salesSummary.fuel.fuel_3.litre"},
			"fuel4":     bson.M{"$sum": "$salesSummary.fuel.fuel_4.litre"},
			"fuel5":     bson.M{"$sum": "$salesSummary.fuel.fuel_5.litre"},
			"fuel6":     bson.M{"$sum": "$salesSummary.fuel.fuel_6.litre"},
			"fuelCosts": bson.M{"$last": "$fuelCosts"},
		},
	}

	project := bson.M{
		"$project": bson.M{
			"recordDate": "$_id.recordDate",
			"stationID":  "$_id.stationID",
			"fuel1":      1,
			"fuel2":      1,
			"fuel3":      1,
			"fuel4":      1,
			"fuel5":      1,
			"fuel6":      1,
			"fuelCosts":  1,
		},
	}

	sort := bson.M{
		"$sort": bson.M{"_id.recordDate": 1},
	}

	pipe := col.Pipe([]bson.M{match, group, project, sort})
	pipe.All(&ss)

	return ss, err
}

func (db *DB) persistFuelSales(ss []*model.StationSales) (ts int64, err error) {

	s := db.getFreshSession()
	defer s.Close()
	col := s.DB(DBSales).C(colFSImport)

	ts = time.Now().Unix()

	for _, elem := range ss {
		fuelSplit := elem.Fuel2 / 2
		rdte, _ := strconv.Atoi(elem.RecordDate.Format(timeShortForm))
		fs := &model.FuelSales{
			NL:   (elem.Fuel1 + fuelSplit),
			SNL:  (elem.Fuel3 + fuelSplit),
			DSL:  elem.Fuel4,
			CDSL: elem.Fuel5,
			PROP: elem.Fuel6,
		}
		fsums := &model.FuelSums{
			Fuel1: elem.Fuel1,
			Fuel2: elem.Fuel2,
			Fuel3: elem.Fuel3,
			Fuel4: elem.Fuel4,
			Fuel5: elem.Fuel5,
			Fuel6: elem.Fuel6,
		}
		fsi := &model.FuelSalesImport{
			FuelCosts:  elem.FuelCosts,
			FuelSales:  fs,
			FuelSums:   fsums,
			ImportTS:   ts,
			RecordDate: rdte,
			StationID:  elem.StationID,
			Status:     "imported",
		}
		if err := col.Insert(fsi); err != nil {
			return ts, err
		}
	}
	return ts, err
}

func (db *DB) compileFuelSales() (err error) {

	// Get list of station nodes to later match with
	nodes, err := db.fetchStationNodes()

	s := db.getFreshSession()
	defer s.Close()
	colIm := s.DB(DBSales).C(colFSImport)
	colEx := s.DB(DBSales).C(colFSExport)

	for _, station := range nodes {
		// found this explaining the $cond operator: https://github.com/go-mgo/mgo/issues/298
		// also this, thou sure how successful: https://groups.google.com/forum/#!topic/mgo-users/yl0eIb0Wh-c
		// also: https://stackoverflow.com/questions/40259171/mongo-aggregation-query-in-golang-with-mgo-driver
		match := bson.M{
			"$match": bson.M{"stationID": bson.M{"$in": station.Nodes}},
		}

		group := bson.M{
			"$group": bson.M{
				"_id": bson.M{"recordDate": "$recordDate", "importTS": "$importTS"},
				"avgFuelCost": bson.M{
					"$avg": bson.M{
						"$cond": bson.M{
							"if":   bson.M{"$gt": []interface{}{"$fuelCosts.fuel_1", 0}},
							"then": "$fuelCosts.fuel_1",
							"else": nil,
						},
					},
				},
				"NL":   bson.M{"$sum": "$fuelSales.NL"},
				"SNL":  bson.M{"$sum": "$fuelSales.SNL"},
				"DSL":  bson.M{"$sum": "$fuelSales.DSL"},
				"CDSL": bson.M{"$sum": "$fuelSales.CDSL"},
				"PROP": bson.M{"$sum": "$fuelSales.PROP"},
			},
		}

		project := bson.M{
			"$project": bson.M{
				"recordDate":  "$_id.recordDate",
				"stationID":   station.ID,
				"avgFuelCost": 1,
				"importTS":    "$_id.importTS",
				"fuelSales": bson.M{
					"NL":   "$NL",
					"SNL":  "$SNL",
					"DSL":  "$DSL",
					"CDSL": "$CDSL",
					"PROP": "$PROP",
				},
			},
		}

		// Oddly, this is not sorting properly
		sort := bson.M{
			"$sort": bson.M{"_id.recordDate": 1},
		}

		var res []model.FuelSalesExport
		pipe := colIm.Pipe([]bson.M{match, group, project, sort})
		pipe.All(&res)

		for _, ex := range res {
			ex.ID = strconv.Itoa(ex.RecordDate) + "-" + ex.StationID.Hex()
			_, err = colEx.Upsert(bson.M{"_id": ex.ID}, ex)
			if err != nil {
				log.Errorf("Error upserting fuel sale export %s", err)
				break
			}
		}
	}

	return err
}

func (db *DB) removeImportedFuelSales() (err error) {

	s := db.getFreshSession()
	defer s.Close()

	col := s.DB(DBSales).C(colFSImport)
	_, err = col.RemoveAll(nil)
	// fmt.Printf("RemoveAll info: %+v\n", info.Removed)

	return err
}

// ==================== Propane methods ==================================== //

func (db *DB) fetchPropaneSales(req *model.Request) (ps []*model.PropaneSale, err error) {

	s := db.getFreshSession()
	defer s.Close()

	col := s.DB(DBSales).C(colFuelSales)

	match := bson.M{
		"$match": bson.M{
			"stationID":  bson.ObjectIdHex(config.PropaneStationID),
			"recordDate": bson.M{"$gte": req.DateStart, "$lte": req.DateEnd},
			"gradeID":    config.PropaneGradeID,
		},
	}

	sort := bson.M{"$sort": bson.M{"recordDate": -1}}

	group := bson.M{
		"$group": bson.M{
			"_id":    bson.M{"recordDate": "$recordDate", "dispenserID": "$dispenserID"},
			"litres": bson.M{"$sum": "$litres.net"},
		},
	}

	project := bson.M{
		"$project": bson.M{
			"recordDate":  "$_id.recordDate",
			"dispenserID": "$_id.dispenserID",
			"litres":      1,
		},
	}

	pipe := col.Pipe([]bson.M{match, sort, group, project})
	pipe.All(&ps)
	return ps, err
}

func (db *DB) persistPropaneSales(ps []*model.PropaneSale) (ts int64, err error) {

	s := db.getFreshSession()
	defer s.Close()
	col := s.DB(DBSales).C(colPSExport)

	ts = time.Now().Unix()

	for _, elem := range ps {
		rdte, _ := strconv.Atoi(elem.RecordDate.Format(timeShortForm))
		psi := &model.PropaneSaleExport{
			ImportTS:   ts,
			Litres:     elem.Litres,
			RecordDate: rdte,
			TankID:     config.PropaneTankLookup(elem.DispenserID.Hex()),
		}
		if err := col.Insert(psi); err != nil {
			return ts, err
		}
	}

	return ts, err
}

// ==================== Fuel & Propane methods ============================= //

func (db *DB) fetchStationNodes() (nodes []*model.StationNodes, err error) {

	s := db.getFreshSession()
	defer s.Close()

	col := s.DB(DBSales).C(colStationNodes)
	col.Find(bson.M{}).All(&nodes)

	return nodes, err
}

func (db *DB) createImportLog(req *model.Request, ts int64) (err error) {

	s := db.getFreshSession()
	defer s.Close()
	col := s.DB(DBSales).C(colImportLog)

	dteSt, _ := strconv.Atoi(req.DateStart.Format(timeShortForm))
	dteEd, _ := strconv.Atoi(req.DateEnd.Format(timeShortForm))

	ilog := &model.ImportLog{
		DateFrom:   dteSt,
		DateTo:     dteEd,
		ImportTS:   ts,
		ImportType: req.ExportType,
	}
	err = col.Insert(ilog)

	return err
}

// ==================== DB Helper methods ==================== //

// Close method
func (db *DB) Close() {
	db.session.Close()
}

func (db *DB) getFreshSession() *mgo.Session {
	return db.session.Copy()
}
