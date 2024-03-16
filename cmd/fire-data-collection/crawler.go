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

	"github.com/Dave-Lopper/airq-backend/internals/types"
)

var apiKey string = os.Getenv("FIRMS_API_KEY")
var baseUrl string = fmt.Sprintf("https://firms.modaps.eosdis.nasa.gov/api/area/csv/%v/MODIS_NRT/88.594382881345,-12.375471986864,152.93774260658,31.291864723093", apiKey)


func CrawlApi(days string, startDate string) ([]types.Fire, error) {
	fullyQualifiedUrl := fmt.Sprintf("%v/%v/%v", baseUrl, days, startDate)
	response, err := http.Get(fullyQualifiedUrl)

	if err != nil {
		log.Fatalf("Error while calling api: %v", err)
		return nil, err
	}
	defer response.Body.Close()

	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		log.Fatalf("Error reading api responde body: %v", readErr)
		return nil, err
	}

	r := csv.NewReader(strings.NewReader(string(body)))

	fires := []types.Fire{}

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

		fires = append(fires, types.Fire{
			Lat: convertRecordVal[float64]("lat", record[0]),
			Lng: convertRecordVal[float64]("lng", record[1]),
			Brightness: convertRecordVal[float64]("brightness", record[2]),
			Scan: convertRecordVal[float64]("scan", record[3]),
			Track: convertRecordVal[float64]("track", record[4]),
			AcqDate: record[5],
			AcqTime: record[6],
			Satellite: record[7],
			Instrument: record[8],
			Confidence: convertRecordVal[int64]("confidence", record[9]),
			Version: record[10],
			BrightT31: convertRecordVal[float64]("brightT31", record[11]),
			Frp: convertRecordVal[float64]("frp", record[12]),
			DayNight: record[13],
		})
	}
	return fires, nil
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