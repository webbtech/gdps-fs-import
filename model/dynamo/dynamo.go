package dynamo

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/pulpfree/gales-fuelsale-export/config"
	"github.com/pulpfree/gales-fuelsale-export/model"
)

// Dynamo struct
type Dynamo struct {
	config *config.Dynamo
	db     *dynamodb.DynamoDB
}

// NewDB connection function
func NewDB(cfg *config.Dynamo) (*Dynamo, error) {

	var err error

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
	})
	if err != nil {
		return nil, err
	}
	svc := dynamodb.New(sess)

	return &Dynamo{
		config: cfg,
		db:     svc,
	}, err
}

// CreateFuelSalesRecords method
func (d *Dynamo) CreateFuelSalesRecords(sales []*model.FuelSalesExport, res *model.DnImportRes) (err error) {

	stations, err := d.fetchStations()

	for _, sale := range sales {

		stationRef := sale.StationID.Hex()
		stationID := stations[stationRef].ID

		fuelSales := &model.FuelSales{
			NL:   sale.FuelSales.NL,
			SNL:  sale.FuelSales.SNL,
			DSL:  sale.FuelSales.DSL,
			CDSL: sale.FuelSales.CDSL,
			PROP: sale.FuelSales.PROP,
		}
		item := model.DnFuelSales{
			AvgFuelCost: sale.AvgFuelCost,
			Date:        sale.RecordDate,
			ImportTS:    sale.ImportTS,
			Sales:       fuelSales,
			StationID:   stationID,
			YearWeek:    setYearWeek(sale.RecordDate),
		}

		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			log.Errorf("Error marshalling map: %s", err)
			return err
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(FuelSale),
		}
		_, err = d.db.PutItem(input)
		if err != nil {
			log.Errorf("Error calling PutItem: %s", err)
			return err
		}

		err = d.createFuelPriceRecord(item)
		if err != nil {
			log.Errorf("Error calling createFuelPriceRecord: %s", err)
			return err
		}
	}

	err = d.createImportLog(res)
	if err != nil {
		log.Errorf("Error calling createImportLog: %s", err)
		return err
	}

	return err
}

// CreatePropaneSalesRecords method
func (d *Dynamo) CreatePropaneSalesRecords(sales []*model.PropaneSaleExport, res *model.DnImportRes) (err error) {

	for _, sale := range sales {
		item := model.DnPropaneSales{
			Date:     sale.RecordDate,
			ImportTS: sale.ImportTS,
			Sales:    sale.Litres,
			TankID:   sale.TankID,
			Year:     setYear(sale.RecordDate),
			YearWeek: setYearWeek(sale.RecordDate),
		}

		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			log.Errorf("Error marshalling map: %s", err)
			return err
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(PropaneSale),
		}
		_, err = d.db.PutItem(input)
		if err != nil {
			log.Errorf("Error calling PutItem: %s", err)
			return err
		}
	}

	err = d.createImportLog(res)
	if err != nil {
		log.Errorf("Error calling createImportLog: %s", err)
		return err
	}

	return err
}

// fetchStations method
func (d *Dynamo) fetchStations() (stationMap map[string]*model.DnStation, err error) {

	proj := expression.NamesList(expression.Name("ID"), expression.Name("Name"), expression.Name("RefStation"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		log.Errorf("Error building expression: %s", err)
		return stationMap, err
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames: expr.Names(),
		FilterExpression:         expr.Filter(),
		ProjectionExpression:     expr.Projection(),
		TableName:                aws.String(Station),
	}

	result, err := d.db.Scan(params)
	if err != nil {
		log.Errorf("Dynamo query API call failed: %s", err)
		return stationMap, err
	}

	stationMap = make(map[string]*model.DnStation)

	for _, i := range result.Items {
		item := &model.DnStation{}
		err := dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Errorf("Error unmarshalling: %s", err)
			return stationMap, err
		}
		stationMap[item.RefStation] = item
	}

	return stationMap, err
}

// createImportLog method
func (d *Dynamo) createImportLog(res *model.DnImportRes) (err error) {

	av, err := dynamodbattribute.MarshalMap(res)
	if err != nil {
		log.Errorf("Error marshalling map: %s", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(ImportLog),
	}
	_, err = d.db.PutItem(input)
	if err != nil {
		log.Errorf("Error calling PutItem: %s", err)
		return err
	}

	return err
}

// createFuelPriceRecord
func (d *Dynamo) createFuelPriceRecord(fs model.DnFuelSales) (err error) {

	item := model.DnFuelPrice{
		Date:      fs.Date,
		Price:     fs.AvgFuelCost,
		StationID: fs.StationID,
		YearWeek:  fs.YearWeek,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Errorf("Error marshalling map: %s", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(FuelPrice),
	}
	_, err = d.db.PutItem(input)
	if err != nil {
		log.Errorf("Error calling PutItem: %s", err)
		return err
	}

	return err
}

// setYearWeek extracts a yearweek YYYYWW integer from the provided date
//
// As golang uses the ISO week with Sunday being the last day of the week,
// this function attempts to use the US/Canada/Australia method with Sunday
// as the first day of the week
//
// It's possible this could create problems at some point in the calendar
func setYearWeek(date int) (yearWeek int) {
	t, _ := time.Parse("20060102", strconv.Itoa(date))
	yr, wk := t.ISOWeek()
	if t.Weekday() == 0 {
		wk++
	}
	yearWeek, _ = strconv.Atoi(fmt.Sprintf("%d%d", yr, wk))

	return yearWeek
}

func setYear(date int) (year int) {
	t, _ := time.Parse("20060102", strconv.Itoa(date))

	return t.Year()
}
