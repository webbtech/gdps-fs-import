package mongo

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/pulpfree/gsales-fs-export/config"
	"github.com/pulpfree/gsales-fs-export/model"
	"github.com/pulpfree/gsales-fs-export/validators"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dateEnd    = "2020-07-15"
	dateStart  = "2020-07-01"
	defaultsFP = "../../config/defaults.yaml"
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

	// Set client options
	clientOptions := options.Client().ApplyURI(s.cfg.GetMongoConnectURL())

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Printf("Error connecting to db: %s", err)
		return
	}

	s.db = &MDB{
		client: client,
		dbName: s.cfg.MongoDBName,
		db:     client.Database(s.cfg.MongoDBName),
	}

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
		ExportType: "fuel",
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

// TestNewDB method
func (s *IntegSuite) TestNewDB() {
	_, err := NewDB(s.cfg.GetMongoConnectURL(), s.cfg.MongoDBName)
	s.NoError(err)
}

// ===================== Un-exported Functions ============================================ //

// TestfetchFuelSales method
func (s *IntegSuite) TestfetchFuelSales() {
	docs, err := s.db.fetchFuelSales(s.fuelReq)
	s.NoError(err)
	s.True(len(docs) > 100)
}

// TestfetchPropaneSales method
func (s *IntegSuite) TestfetchPropaneSales() {
	docs, err := s.db.fetchPropaneSales(s.fuelReq)
	s.NoError(err)
	// s.True(len(docs) > 100)
	fmt.Printf("docs: %+v\n", docs)
}
