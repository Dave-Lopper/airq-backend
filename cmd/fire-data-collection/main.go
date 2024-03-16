package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Dave-Lopper/airq-backend/internals/db"
)

func main() {
	days := flag.Int("days", 7, "Number of days from provided date")
	startDateStr := flag.String("start_date", "", "Date to start crawling from (format YYYY-MM-DD)")
	flag.Parse()

	dateLayout := "2006-01-02"
	var startDate time.Time
	now := time.Now()

	if *days > 10 {
		log.Fatalf("Can't parse more than 10 days worth of data at a time")
		return
	}

	if *startDateStr != "" {
		var err error
		startDate, err = time.Parse(dateLayout, *startDateStr)
		if err != nil {
			log.Fatalf("Invalid date format %v, expected YYYY-MM-DD", *startDateStr)
			return
		}

		if startDate.After(now) {
			log.Fatalf("Provided date %v is in the future", *startDateStr)
			return
		}
	} else {
		startDate = now
	}

	fires, err := CrawlApi(strconv.Itoa(*days), startDate.Format(dateLayout))
	if err != nil {
		log.Fatalf("Error crawling API: %v", err)
		panic(err)
	}

	goquDB := db.GetConnection()
	exists, err := db.TableExists(goquDB, "t_fires")
	if err != nil {
		log.Fatalf("Error while checking if table exists: %v", err)
		panic(err)
	}

	if exists {
		fmt.Println("Table t_fires exists")
	} else {
		fmt.Println("Table t_fires doesn't exist, creating it")
		err := db.CreateTable(goquDB)
		if err != nil {
			log.Fatalf("Error while creating t_fires table: %v", err)
			panic(err)
		}
	}


	insertErr := db.InsertFires(goquDB, fires)
	if insertErr != nil {
		panic(insertErr)
	}
}
