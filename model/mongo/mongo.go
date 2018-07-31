package mongo

import (
	"github.com/pulpfree/gales-fuelsale-export/model"
	mgo "gopkg.in/mgo.v2"
)

// DB struct
type DB struct {
	session *mgo.Session
}

// DB Constants
const (
	DBSales  = "gales-sales"
	colSales = "sales"
)

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
func (db *DB) CreateFuelSales(req *model.Request) error {

	return nil
}

// Close method
func (db *DB) Close() {
	db.session.Close()
}

// Helper methods

func (db *DB) getFreshSession() *mgo.Session {
	return db.session.Copy()
}
