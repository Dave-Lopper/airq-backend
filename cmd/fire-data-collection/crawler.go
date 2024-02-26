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

		fires = append(fires, Fire{
			record[0],
			record[1],
			convertRecordVal[float64]("brightness", record[2]),
			convertRecordVal[float64]("scan", record[3]),
			convertRecordVal[float64]("track", record[4]),
			record[5],
			record[6],
			record[7],
			record[8],
			convertRecordVal[int64]("confidence", record[9]),
			record[10],
			convertRecordVal[float64]("brightT31", record[11]),
			convertRecordVal[float64]("frp", record[12]),
			record[13],
		})
	}

	fmt.Printf("Fires: %v", fires)
}

func convertRecordVal[T int64 | float64](name string, recordVal string) T {
	if strings.ContainsAny(recordVal, ".") {
		convertedVal, err := strconv.ParseFloat(recordVal, 64)
		if err != nil {
			fmt.Printf("Error converting %v to float: %v", name, err)
			panic(err)
		} 
		return T(convertedVal)
	} else {
		convertedVal, err := strconv.ParseInt(recordVal, 10, 64)
		if err != nil {
			fmt.Printf("Error converting %v to int: %v", name, err)
			panic(err)
		}
		return T(convertedVal)
	}
}