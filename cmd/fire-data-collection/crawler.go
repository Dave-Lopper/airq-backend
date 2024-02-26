package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var apiKey string = os.Getenv("FIRMS_API_KEY")
var baseUrl string = fmt.Sprintf("https://firms.modaps.eosdis.nasa.gov/api/area/csv/%v/MODIS_NRT/88.594382881345,-12.375471986864,152.93774260658,31.291864723093", apiKey)

type Fire struct {
	lat        string
	lng        string
	brightness float64
	scan       float64
	track      float64
	acqDate    string
	acqTime    string
	satellite  string
	instrument string
	confidence int64
	version    string
	brightT31  float64
	frp        float64
	dayNight   string
}

func CrawlApi(days string, startDate string) {
	fullyQualifiedUrl := fmt.Sprintf("%v/%v/%v", baseUrl, days, startDate)
	fmt.Printf("Full URL: %v", fullyQualifiedUrl)
	response, err := http.Get(fullyQualifiedUrl)

	if err != nil {
		log.Fatalf("Error while calling api: %v", err)
		return
	}
	defer response.Body.Close()

	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		log.Fatalf("Error reading api responde body: %v", readErr)
		return
	}

	r := csv.NewReader(strings.NewReader(string(body)))

	fires := []Fire{}

	i := 0

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		i += 1
		if i == 1 {
			continue
		}

		if err != nil {
			log.Fatal(err)
		}

		brightness, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			fmt.Printf("Error converting brightness to float: %v", err)
			return
		}

		scan, scanErr := strconv.ParseFloat(record[3], 64)
		if scanErr != nil {
			fmt.Printf("Error converting scan to float: %v", err)
			return
		}

		track, trackErr := strconv.ParseFloat(record[4], 64)
		if trackErr != nil {
			fmt.Printf("Error converting track to float: %v", err)
			return
		}

		confidence, confidenceErr := strconv.ParseInt(record[9], 10, 64)
		if confidenceErr != nil {
			fmt.Printf("Error converting confidence to int: %v", err)
			return
		}

		brightT31, brightT31Err := strconv.ParseFloat(record[11], 64)
		if brightT31Err != nil {
			fmt.Printf("Error converting brightT31 to float: %v", err)
			return
		}

		frp, frpErr := strconv.ParseFloat(record[12], 64)
		if frpErr != nil {
			fmt.Printf("Error converting frp to float: %v", err)
			return
		}

		fire := Fire{
			record[0],
			record[1],
			brightness,
			scan,
			track,
			record[5],
			record[6],
			record[7],
			record[8],
			confidence,
			record[10],
			brightT31,
			frp,
			record[13],
		}
		fmt.Printf("%v", fire)
		fires = append(fires, fire)
	}

	fmt.Printf("Fires: %v", fires)
}
