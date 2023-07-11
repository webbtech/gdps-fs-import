package mongo

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pulpfree/gsales-fs-export/config"
	"github.com/pulpfree/gsales-fs-export/model"
	"github.com/pulpfree/gsales-fs-export/validators"
	"github.com/stretchr/testify/suite"
)

const (
	// dateEnd    = "2020-07-05"
	// dateStart  = "2020-07-01"
	dateEnd    = "2023-06-30"
	dateStart  = "2023-06-06"
	defaultsFP = "../../config/defaults.yml"
	timeForm   = "2006-01-02"
)

// IntegSuite struct
type IntegSuite struct {
	cfg     *config.Config
	db      *MDB
	fuelReq *model.Request
	propReq *model.Request
	suite.Suite
}

// SetupTest method
func (s *IntegSuite) SetupTest() {
	// setup config
	os.Setenv("Stage", "test")
	s.cfg = &config.Config{DefaultsFilePath: defaultsFP}
	err := s.cfg.Load()
	if err != nil {
		fmt.Printf("Error in loading config: %s", err)
		return
	}
	s.NoError(err)

	s.db, err = NewDB(s.cfg.GetMongoConnectURL(), s.cfg.MongoDBName)
	if err != nil {
		fmt.Printf("Error connecting to db: %s", err)
		return
	}

	// create fuel and propane requests
	fuelTestVars := &model.RequestInput{
		DateStart:  dateStart,
		DateEnd:    dateEnd,
		ExportType: "fuel",
	}
	s.fuelReq, err = validators.RequestVars(fuelTestVars)
	if err != nil {
		fmt.Printf("Error validating fuel request: %s", err)
		return
	}

	propTestVars := &model.RequestInput{
		DateStart:  dateStart,
		DateEnd:    dateEnd,
		ExportType: "propane",
	}
	s.propReq, err = validators.RequestVars(propTestVars)
	if err != nil {
		fmt.Printf("Error validating propane request: %s", err)
		return
	}
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}

// ===================== Exported Functions ================================================ //

// TestCreateFuelSales method
func (s *IntegSuite) TestCreateFuelSales() {
	defer s.db.Close()

	err := s.db.CreateFuelSales(s.fuelReq)
	s.NoError(err)
}

// TestCreatePropaneSales method
func (s *IntegSuite) TestCreatePropaneSales() {
	defer s.db.Close()

	err := s.db.CreatePropaneSales(s.propReq)
	s.NoError(err)
}

// TestFetchExportedFuelSales method
func (s *IntegSuite) TestFetchExportedFuelSales() {
	defer s.db.Close()

	docs, err := s.db.FetchExportedFuelSales(s.fuelReq)
	s.NoError(err)
	s.True(len(docs) > 10)
}

// TestFetchExportedPropaneSales method
func (s *IntegSuite) TestFetchExportedPropaneSales() {
	defer s.db.Close()

	docs, err := s.db.FetchExportedPropaneSales(s.propReq)
	s.NoError(err)
	s.True(len(docs) > 2)
}

// ===================== Un-exported Functions ============================================ //

// TestfetchFuelSales method
func (s *IntegSuite) TestfetchFuelSales() {
	defer s.db.Close()

	docs, err := s.db.fetchFuelSales(s.fuelReq)
	s.NoError(err)
	s.True(len(docs) > 10)
}

// TestfetchPropaneSales method
func (s *IntegSuite) TestfetchPropaneSales() {
	defer s.db.Close()

	docs, err := s.db.fetchPropaneSales(s.fuelReq)
	s.NoError(err)
	s.True(len(docs) > 2)
}

// TestpersistFuelSales method
func (s *IntegSuite) TestpersistFuelSales() {
	defer s.db.Close()

	docs, err := s.db.fetchFuelSales(s.fuelReq)
	s.NoError(err)

	ts1 := time.Now().Unix()
	ts2, err := s.db.persistFuelSales(docs)
	s.NoError(err)
	s.True(ts1 == ts2)
}

// TestfetchStationNodes method
func (s *IntegSuite) TestfetchStationNodes() {
	defer s.db.Close()

	nodes, err := s.db.fetchStationNodes()
	s.NoError(err)
	s.True(len(nodes) > 10)
}

// TestremoveImportedFuelSales method
// this should be run last
func (s *IntegSuite) TestremoveImportedFuelSales() {
	defer s.db.Close()

	res, err := s.db.removeImportedFuelSales()
	s.NoError(err)
	s.True(res.DeletedCount > 10)
}

// TestcompileFuelSales method
func (s *IntegSuite) TestcompileFuelSales() {
	defer s.db.Close()

	_, err := s.db.removeImportedFuelSales()
	s.NoError(err)

	docs, err := s.db.fetchFuelSales(s.fuelReq)
	s.NoError(err)

	_, err = s.db.persistFuelSales(docs)
	s.NoError(err)

	err = s.db.compileFuelSales()
	s.NoError(err)

	_, err = s.db.removeImportedFuelSales()
	s.NoError(err)
}

// TestcreateImportLog method
func (s *IntegSuite) TestcreateImportLog() {
	defer s.db.Close()

	ts := time.Now().Unix()
	_, err := s.db.createImportLog(s.fuelReq, ts)
	s.NoError(err)
}

// TestpersistPropaneSales method
func (s *IntegSuite) TestpersistPropaneSales() {
	defer s.db.Close()

	docs, err := s.db.fetchPropaneSales(s.fuelReq)
	s.NoError(err)
	fmt.Printf("docs: %+v\n", docs[0])

	ts, err := s.db.persistPropaneSales(docs)
	fmt.Printf("ts: %+v\n", ts)
}
