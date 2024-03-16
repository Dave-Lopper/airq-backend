package db

import (
	"fmt"
	"log"

	"github.com/Dave-Lopper/airq-backend/internals/types"
	"github.com/doug-martin/goqu/v9"
)


func TableExists(db *goqu.Database, tableName string) (bool, error) {
	query, _, err := db.
		From(goqu.I("information_schema.tables")).
		Select(goqu.C("table_name")).
		Where(
			goqu.Ex{"table_name": tableName},
		).
		Limit(1).ToSQL()

	if err != nil {
		return false, err
	}

	var result string
	err = db.QueryRow(query).Scan(&result)
	
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		} else {
			log.Fatalf("Error while checking if table exists: %v", err)
			return false, err
		}
	}
	return true, nil
}


func CreateTable(db *goqu.Database) error {
	_, err := db.Exec(`
		CREATE TABLE t_fires (
			location geometry(Point, 4326),
			brightness FLOAT,
			scan FLOAT,
			track FLOAT,
			acq_date timestamp,
			satellite VARCHAR(255),
			instrument VARCHAR(255),
			confidence INT,
			version VARCHAR(255),
			bright_t31 FLOAT,
			frp FLOAT,
			day_night VARCHAR(255)
		)
	`)
	return err
}


func InsertFires(db *goqu.Database, fires []types.Fire) error {
	var records []goqu.Record

	for _, fire := range fires {
		record := goqu.Record{
			"location": goqu.Func("ST_SetSRID", goqu.Func("ST_MakePoint", fire.Lat, fire.Lng), 4326),
			"brightness": fire.Brightness,
			"scan": fire.Scan,
			"track": fire.Track,
			"acq_date": fire.AcqDate,
			"satellite": fire.Satellite,
			"instrument": fire.Instrument,
			"confidence": fire.Confidence,
			"version": fire.Version,
			"bright_t31": fire.BrightT31,
			"frp": fire.Frp,
			"day_night": fire.DayNight,
		}
		records = append(records, record)
	}

	var interfaceRecords []interface{}
	for _, record := range records {
		interfaceRecords = append(interfaceRecords, record)
	}

	ds := db.Insert("t_fires").Rows(interfaceRecords...)
	insertSQL, args, _ := ds.ToSQL()
	fmt.Println(insertSQL)
	fmt.Println(args)
	_, insertErr := db.Exec(insertSQL, args...)

	if insertErr != nil {
		return insertErr
	}
	return nil
}