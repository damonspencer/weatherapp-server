package weatherdata

/*
{
"Country" : "Some Country"
"Threedigitcc":"SMC", //somecountry's three digit country code
"City":"Somecity",
"Lat":123.456,
"Lon":789.012,
"Weather":"cloudy",
"Temperature":"314",  //in kelvin
"Windspeed":5.1, //mph
"Winddir":150, //in degrees
"Percentcloudy":74, //percent
"Percentrain":"0", //percent
"time" : "21:34" //24 hour time format
}
*/
type Weatherdata struct {
	Country       string
	Threedigitcc  string
	City          string
	Lat           float64
	Lon           float64
	Weather       string
	Temperature   int
	Windspeed     float64
	Winddir       float64
	Percentcloudy int
	Percentrain   int
	Hour          int
	Min           int
}
