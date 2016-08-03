package tests

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

//gets data from testdata and updates the server every minute
//for now, the fake weather data is in the format:
/*
{
"Country":"Somecountry",
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

//dynamically generate and post the fake weather data

var country string = "Magnificent Glorious Country"
var cities []string = []string{"Foo", "Bar", "Fgh", "somecity", "whatami"}
var lats []float64 = []float64{123.456, 654.321, 732.345, 908.243, 675.346, 567.658}
var lons []float64 = []float64{864.468, 324.987, 576.243, 456.899, 134.456, 466.787}

//made with the world's highest quality rng:
//smashing the number keys on your keyboard
var weathertypes []string = []string{"clear",
	"cloudy", "rainy", "snowy", "sunny"}

var wdslice []Weatherdata

//this should be sufficent for testing

func Fakewdhttpserver() {

	t := time.Now() //get time
	r, err := strconv.Atoi(t.Format("20060102150405"))
	if err != nil {
		fmt.Println(err)
		//TODO better error recovery
	}
	rand.Seed(int64(r)) //seed the rng
	populatewd()
	//initially populate the weatherdata
	//slice so we don't have null indexes
	http.HandleFunc("/fakewd", handler)
	http.ListenAndServe(":5555", nil)
	//a bad, but working handler function because
	//the better ones I made didn't work right
}

//dynamically generate weather data each time you refresh the page
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, getfakewd(), r.URL.Path[1:])
}

func getfakewd() string {
	newwdslice := make([]Weatherdata, 5)
	//create a slice of weatherdata structs
	for i := range newwdslice {
		weathertype := ""
		percentcloudy := 0
		percentrain := 0
		temp := 0
		windspeed := 0.0
		winddir := 0.0
		hour := wdslice[i].Hour
		min := wdslice[i].Min + rand.Intn(20)
		if min >= 60 {
			min = min - 60
			hour++
		}
		if rand.Intn(2) == 1 || ((hour >= 18 || hour <= 6) && wdslice[i].Weather == "sunny") {
			//only a 50% chance for the weather to change,
			//unless its night and the old weather says its sunny,
			//in which case the weather changes
			if hour >= 18 || hour <= 6 { //its night, it can't be sunny
				weathertype = weathertypes[rand.Intn(len(weathertypes)-1)]
			} else { //its day
				weathertype = weathertypes[rand.Intn(len(weathertypes))]
			}
			//if its clear or sunny it can't be cloudy or rainy
			if weathertype != "clear" && weathertype != "sunny" {
				percentcloudy = rand.Intn(100)
				percentrain = rand.Intn(100)
			}
			temp = rand.Intn(100) + 220
			windspeed = rand.Float64() * 100 //arbitrary max speed
			winddir = rand.Float64() * 360
		} else { //no change, so use the old weather data
			weathertype = wdslice[i].Weather
			percentcloudy = wdslice[i].Percentcloudy
			percentrain = wdslice[i].Percentrain
			temp = wdslice[i].Temperature
			windspeed = wdslice[i].Windspeed
			winddir = wdslice[i].Winddir
		}
		newwdslice[i] = Weatherdata{country, cities[i], lats[i],
			lons[i], weathertype, temp, windspeed, winddir,
			percentcloudy, percentrain, hour, min}
		//create an instance of weatherdata in the slice
	}
	wdjson, err := json.Marshal(newwdslice)
	//jsonify
	wdjsonstr := string(wdjson)
	if err != nil {
		fmt.Println(err)
		//TODO better error recovery
	}
	return wdjsonstr
}
func populatewd() {
	wdslice = make([]Weatherdata, 5)
	for i := range wdslice {
		//create a slice of weatherdata structs
		hour := rand.Intn(24)
		min := rand.Intn(60)
		weathertype := ""
		percentcloudy := 0
		percentrain := 0
		if hour >= 18 || hour <= 6 { //its night, it can't be sunny
			weathertype = weathertypes[rand.Intn(len(weathertypes)-1)]
		} else { //its day
			weathertype = weathertypes[rand.Intn(len(weathertypes))]
		}
		//if its clear or sunny it can't be cloudy or rainy
		if weathertype != "clear" && weathertype != "sunny" {
			percentcloudy = rand.Intn(100)
			percentrain = rand.Intn(100)
		}
		temp := rand.Intn(100) + 220
		windspeed := rand.Float64() * 100 //arbitrary max speed
		winddir := rand.Float64() * 360

		wdslice[i] = Weatherdata{country, cities[i], lats[i], lons[i],
			weathertype, temp, windspeed, winddir, percentcloudy,
			percentrain, hour, min}
		//create an instance of weatherdata
	}
}
