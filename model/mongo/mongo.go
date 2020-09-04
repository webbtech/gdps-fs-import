package mongo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/pulpfree/gsales-fs-export/config"
	"github.com/pulpfree/gsales-fs-export/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MDB struct
type MDB struct {
	client *mongo.Client
	dbName string
	db     *mongo.Database
}

// DB and collections Constants
const (
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
func NewDB(connection string, dbNm string) (*MDB, error) {

	clientOptions := options.Client().ApplyURI(connection)
	err := clientOptions.Validate()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	return &MDB{
		client: client,
		dbName: dbNm,
		db:     client.Database(dbNm),
	}, err
}

// CreateFuelSales function
func (db *MDB) CreateFuelSales(req *model.Request) (err error) {

	sales, err := db.fetchFuelSales(req)
	if err != nil {
		return err
	}

	ts, err := db.persistFuelSales(sales)
	if err != nil {
		return err
	}

	_, err = db.createImportLog(req, ts)
	if err != nil {
		return err
	}

	err = db.compileFuelSales()
	if err != nil {
		return err
	}

	_, err = db.removeImportedFuelSales()
	if err != nil {
		return err
	}

	return err
}

// CreatePropaneSales function
func (db *MDB) CreatePropaneSales(req *model.Request) (err error) {

	sales, err := db.fetchPropaneSales(req)
	if err != nil {
		return err
	}

	ts, err := db.persistPropaneSales(sales)
	if err != nil {
		return err
	}

	_, err = db.createImportLog(req, ts)
	if err != nil {
		return err
	}

	return err
}

// FetchExportedFuelSales method
func (db *MDB) FetchExportedFuelSales(req *model.Request) (docs []*model.FuelSalesExport, err error) {

	// fetch previously exported records by date range
	col := db.db.Collection(colFSExport)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stDte, _ := strconv.Atoi(req.DateStart.Format(timeShortForm))
	enDte, _ := strconv.Atoi(req.DateEnd.Format(timeShortForm))

	filter := bson.D{
		primitive.E{
			Key: "recordDate",
			Value: bson.D{
				primitive.E{
					Key:   "$gte",
					Value: stDte,
				},
				primitive.E{
					Key:   "$lte",
					Value: enDte,
				},
			},
		},
	}
	cur, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, err
}

// FetchExportedPropaneSales method
func (db *MDB) FetchExportedPropaneSales(req *model.Request) (docs []*model.PropaneSaleExport, err error) {

	// fetch previously exported records by date range
	col := db.db.Collection(colPSExport)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stDte, _ := strconv.Atoi(req.DateStart.Format(timeShortForm))
	enDte, _ := strconv.Atoi(req.DateEnd.Format(timeShortForm))

	filter := bson.D{
		primitive.E{
			Key: "recordDate",
			Value: bson.D{
				primitive.E{
					Key:   "$gte",
					Value: stDte,
				},
				primitive.E{
					Key:   "$lte",
					Value: enDte,
				},
			},
		},
	}
	cur, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, err
}

// ==================== FuelSales methods ==================== //

func (db *MDB) fetchFuelSales(req *model.Request) (docs []model.StationSales, err error) {

	col := db.db.Collection(colSales)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{
			primitive.E{
				Key: "$match",
				Value: bson.D{
					primitive.E{
						Key: "recordDate",
						Value: bson.D{
							primitive.E{
								Key:   "$gte",
								Value: req.DateStart,
							},
							primitive.E{
								Key:   "$lte",
								Value: req.DateEnd,
							},
						},
					},
				},
			},
		},
		{
			primitive.E{
				Key: "$group",
				Value: bson.D{
					primitive.E{
						Key: "_id",
						Value: bson.D{
							primitive.E{
								Key:   "recordDate",
								Value: "$recordDate",
							},
							primitive.E{
								Key:   "stationID",
								Value: "$stationID",
							},
						},
					},
					primitive.E{
						Key: "fuel1",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_1.litre",
							},
						},
					},
					primitive.E{
						Key: "fuel2",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_2.litre",
							},
						},
					},
					primitive.E{
						Key: "fuel3",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_3.litre",
							},
						},
					},
					primitive.E{
						Key: "fuel4",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_4.litre",
							},
						},
					},
					primitive.E{
						Key: "fuel5",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_5.litre",
							},
						},
					},
					primitive.E{
						Key: "fuel6",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_6.litre",
							},
						},
					},
					primitive.E{
						Key: "fuelCosts",
						Value: bson.D{
							primitive.E{
								Key:   "$last",
								Value: "$fuelCosts",
							},
						},
					},
				},
			},
		},
		{
			primitive.E{
				Key: "$project",
				Value: bson.D{
					primitive.E{
						Key:   "recordDate",
						Value: "$_id.recordDate",
					},
					primitive.E{
						Key:   "stationID",
						Value: "$_id.stationID",
					},
					primitive.E{
						Key:   "fuel1",
						Value: 1,
					},
					primitive.E{
						Key:   "fuel2",
						Value: 1,
					},
					primitive.E{
						Key:   "fuel3",
						Value: 1,
					},
					primitive.E{
						Key:   "fuel4",
						Value: 1,
					},
					primitive.E{
						Key:   "fuel5",
						Value: 1,
					},
					primitive.E{
						Key:   "fuel6",
						Value: 1,
					},
					primitive.E{
						Key:   "fuelCosts",
						Value: 1,
					},
				},
			},
		},
		{
			primitive.E{
				Key: "$sort",
				Value: bson.D{
					primitive.E{
						Key:   "_id.recordDate",
						Value: 1,
					},
				},
			},
		},
	}

	cur, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, err
}

func (db *MDB) persistFuelSales(docs []model.StationSales) (ts int64, err error) {

	col := db.db.Collection(colFSImport)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ts = time.Now().Unix()
	for _, elem := range docs {
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

		if _, err := col.InsertOne(ctx, fsi); err != nil {
			return ts, err
		}
	}

	return ts, err
}

func (db *MDB) compileFuelSales() (err error) {
	// Get list of station nodes to later match with
	nodes, err := db.fetchStationNodes()

	colIm := db.db.Collection(colFSImport)
	colEx := db.db.Collection(colFSExport)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	for _, station := range nodes {

		pipeline := mongo.Pipeline{
			{
				primitive.E{
					Key: "$match",
					Value: bson.D{
						primitive.E{
							Key: "stationID",
							Value: bson.D{
								primitive.E{
									Key:   "$in",
									Value: station.Nodes,
								},
							},
						},
					},
				},
			},
			{
				primitive.E{
					Key: "$group",
					Value: bson.D{
						primitive.E{
							Key: "_id",
							Value: bson.D{
								primitive.E{
									Key:   "recordDate",
									Value: "$recordDate",
								},
								primitive.E{
									Key:   "importTS",
									Value: "$importTS",
								},
							},
						},
						primitive.E{
							Key: "avgFuelCost",
							Value: bson.D{
								primitive.E{
									Key: "$avg",
									Value: bson.D{
										primitive.E{
											Key: "$cond",
											Value: bson.D{
												primitive.E{
													Key: "if",
													Value: bson.D{
														primitive.E{
															Key:   "$gt",
															Value: []interface{}{"$fuelCosts.fuel_1", 0},
														},
													},
												},
												primitive.E{
													Key:   "then",
													Value: "$fuelCosts.fuel_1",
												},
												primitive.E{
													Key:   "else",
													Value: nil,
												},
											},
										},
									},
								},
							},
						},
						primitive.E{
							Key: "NL",
							Value: bson.D{
								primitive.E{
									Key:   "$sum",
									Value: "$fuelSales.NL",
								},
							},
						},
						primitive.E{
							Key: "SNL",
							Value: bson.D{
								primitive.E{
									Key:   "$sum",
									Value: "$fuelSales.SNL",
								},
							},
						},
						primitive.E{
							Key: "DSL",
							Value: bson.D{
								primitive.E{
									Key:   "$sum",
									Value: "$fuelSales.DSL",
								},
							},
						},
						primitive.E{
							Key: "CDSL",
							Value: bson.D{
								primitive.E{
									Key:   "$sum",
									Value: "$fuelSales.CDSL",
								},
							},
						},
						primitive.E{
							Key: "PROP",
							Value: bson.D{
								primitive.E{
									Key:   "$sum",
									Value: "$fuelSales.PROP",
								},
							},
						},
					},
				},
			},
			{
				primitive.E{
					Key: "$project",
					Value: bson.D{
						primitive.E{
							Key:   "recordDate",
							Value: "$_id.recordDate",
						},
						primitive.E{
							Key:   "stationID",
							Value: station.ID,
						},
						primitive.E{
							Key:   "avgFuelCost",
							Value: 1,
						},
						primitive.E{
							Key:   "importTS",
							Value: "$_id.importTS",
						},
						primitive.E{
							Key:   "_id",
							Value: 0,
						},
						primitive.E{
							Key: "fuelSales",
							Value: bson.D{
								primitive.E{
									Key:   "NL",
									Value: "$NL",
								},
								primitive.E{
									Key:   "SNL",
									Value: "$SNL",
								},
								primitive.E{
									Key:   "DSL",
									Value: "$DSL",
								},
								primitive.E{
									Key:   "CDSL",
									Value: "$CDSL",
								},
								primitive.E{
									Key:   "PROP",
									Value: "$PROP",
								},
							},
						},
					},
				},
			},
			{
				primitive.E{
					Key: "$sort",
					Value: bson.D{
						primitive.E{
							Key:   "recordDate",
							Value: 1,
						},
					},
				},
			},
		}

		cur, err := colIm.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}
		defer cur.Close(ctx)

		var docs []model.FuelSalesExport
		if err := cur.All(ctx, &docs); err != nil {
			return err
		}

		// now we can insert/update fuel export doc
		opts := options.Update().SetUpsert(true)
		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		defer cancel()

		for _, doc := range docs {
			doc.ID = fmt.Sprintf("%s-%s", strconv.Itoa(doc.RecordDate), doc.StationID.Hex())

			filter := bson.D{
				primitive.E{
					Key:   "_id",
					Value: doc.ID,
				},
			}
			update := bson.D{
				primitive.E{
					Key:   "$set",
					Value: doc,
				},
			}
			_, err := colEx.UpdateOne(ctx, filter, update, opts)
			if err != nil {
				log.Errorf("Error upserting fuel sale export. Error: %s", err)
				break
			}
		}
	}

	return err
}

func (db *MDB) removeImportedFuelSales() (res *mongo.DeleteResult, err error) {

	col := db.db.Collection(colFSImport)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	res, err = col.DeleteMany(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	return res, err
}

// ==================== Propane methods ==================================== //

func (db *MDB) fetchPropaneSales(req *model.Request) (docs []model.PropaneSale, err error) {

	col := db.db.Collection(colFuelSales)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	propStationID, _ := primitive.ObjectIDFromHex(config.PropaneStationID)

	pipeline := mongo.Pipeline{
		{
			primitive.E{
				Key: "$match",
				Value: bson.D{
					primitive.E{
						Key:   "stationID",
						Value: propStationID,
					},
					primitive.E{
						Key: "recordDate",
						Value: bson.D{
							primitive.E{
								Key:   "$gte",
								Value: req.DateStart,
							},
							primitive.E{
								Key:   "$lte",
								Value: req.DateEnd,
							},
						},
					},
					primitive.E{
						Key:   "gradeID",
						Value: config.PropaneGradeID,
					},
				},
			},
		},
		{
			primitive.E{
				Key: "$group",
				Value: bson.D{
					primitive.E{
						Key: "_id",
						Value: bson.D{
							primitive.E{
								Key:   "recordDate",
								Value: "$recordDate",
							},
							primitive.E{
								Key:   "dispenserID",
								Value: "$dispenserID",
							},
						},
					},
					primitive.E{
						Key: "litres",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$litres.net",
							},
						},
					},
				},
			},
		},
		{
			primitive.E{
				Key: "$project",
				Value: bson.D{
					primitive.E{
						Key:   "recordDate",
						Value: "$_id.recordDate",
					},
					primitive.E{
						Key:   "dispenserID",
						Value: "$_id.dispenserID",
					},
					primitive.E{
						Key:   "litres",
						Value: 1,
					},
				},
			},
		},
		{
			primitive.E{
				Key: "$sort",
				Value: bson.D{
					primitive.E{
						Key:   "_id.recordDate",
						Value: 1,
					},
				},
			},
		},
	}

	cur, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, err
}

func (db *MDB) persistPropaneSales(docs []model.PropaneSale) (ts int64, err error) {

	col := db.db.Collection(colPSExport)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ts = time.Now().Unix()

	for _, doc := range docs {
		rdte, _ := strconv.Atoi(doc.RecordDate.Format(timeShortForm))
		psi := &model.PropaneSaleExport{
			ImportTS:   ts,
			Litres:     doc.Litres,
			RecordDate: rdte,
			TankID:     config.PropaneTankLookup(doc.DispenserID.Hex()),
		}

		if _, err := col.InsertOne(ctx, psi); err != nil {
			return ts, err
		}
	}

	return ts, err
}

// ==================== Fuel & Propane methods ============================= //

func (db *MDB) fetchStationNodes() (nodes []model.StationNodes, err error) {

	col := db.db.Collection(colStationNodes)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	cur, err := col.Find(ctx, bson.D{})
	if err != nil {
		return nodes, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &nodes); err != nil {
		return nodes, err
	}

	return nodes, err
}

func (db *MDB) createImportLog(req *model.Request, ts int64) (res *mongo.InsertOneResult, err error) {

	col := db.db.Collection(colImportLog)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dteSt, _ := strconv.Atoi(req.DateStart.Format(timeShortForm))
	dteEd, _ := strconv.Atoi(req.DateEnd.Format(timeShortForm))

	importlog := &model.ImportLog{
		DateFrom:   dteSt,
		DateTo:     dteEd,
		ImportTS:   ts,
		ImportType: req.ExportType,
	}

	res, err = col.InsertOne(ctx, importlog)
	if err != nil {
		return nil, err
	}

	return res, err
}

// ==================== DB Helper methods ==================== //

// Close method
func (db *MDB) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.client.Disconnect(ctx); err != nil {
		panic(err)
	}

	log.Println("MongoDB Disconnected")
}
