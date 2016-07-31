package main

import (
	"fmt"
	"net/http"
	"time"
)

var weathermap map[string]string

func wdclient(args []string) {
	//test := false
	//if contains(args, "test") { test = true }
	//test := true
	url := "" //TODO add location of test http server
	weathermap = make(map[string]string)
	//map of location to the actual weather at that place
	//make sure the socket server is running before taking its data
	locationchannel := make(chan []string)
	for { //run until the program exits
		locationchannel <- locationarray
		locations := <-locationchannel
		//fmt.Println(locations)
		gettestwd(url, locations) //only test mode is implemented
		time.Sleep(time.Second)   //sleep for 1 second
	}

}
func gettestwd(url string, locations []string) { //map[string]string {
	//returns a map of wd to locations
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		//TODO better error recovery
	}
	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)

}
func getdeltamap(oldmap map[string]string, newmap map[string]string) {

}

//for reference only send locations > get wd > send back delta map
