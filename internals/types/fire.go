package types

type Fire struct {
	Lat        float64
	Lng        float64
	Brightness float64
	Scan       float64
	Track      float64
	AcqDate    string
	AcqTime    string
	Satellite  string
	Instrument string
	Confidence int64
	Version    string
	BrightT31  float64
	Frp        float64
	DayNight   string
}