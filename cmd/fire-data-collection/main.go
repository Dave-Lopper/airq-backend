package main

import (
	"flag"
	"log"
	"strconv"
	"time"
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

	CrawlApi(strconv.Itoa(*days), startDate.Format(dateLayout))
}
