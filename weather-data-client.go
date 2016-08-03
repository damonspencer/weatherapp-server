package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"weatherapp-server/tests"
)

var oldweathermap map[string]string
var oldwd []tests.Weatherdata

func wdclient(args []string) {
	//test := false
	//if contains(args, "test") { test = true }
	//test := true
	url := "http://localhost:5555/fakewd"
	time.Sleep(time.Second * 10)
	//wait for the fake weather data server to start up
	//TODO find a better way
	oldweathermap = make(map[string]string)
	//map of locations to the actual weather at that place
	//make sure the socket server is running before taking its data

	//just for now, so that I can test the gettestwd function without clients
	//locationchannel := make(chan []string)
	for { //run until the program exits
		//locationchannel <- locationarray
		//TODO figure out if this line be outside of the loop
		//locations := <-locationchannel
		//fmt.Println(locations)
		///gettestwd(url, locations) //only test mode is implemented
		gettestwd(url)
		time.Sleep(time.Second) //sleep for 1 second
	}

}

/////func gettestwd(url string, locations []string) { //map[string]string {
func gettestwd(url string) { //map[string]string {
	//returns a map of wd to locations

	//we're not using a channel because
	//we're trying to pretend this is an unrelated website
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		//TODO better error recovery
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		//TODO better error recovery
	}
	fmt.Println(body)
	fmt.Println(string(body))
	var newwd []tests.Weatherdata
	err = json.Unmarshal(body[:len(body)-23], &newwd)
	//"%!(EXTRA string=fakewd)" appears at the end, so discard it
	//a bad way of extracting a json
	//from a byte array of a json turned into a
	//string turned back into a byte array...
	//maybe I could just discard everything after "%"?
	//or maybe I could make a better http server
	//so I wouldn't have to do this in the first place...

	fmt.Println(newwd)
}
func getdeltamap(oldmap map[string]string, newmap map[string]string) {

}

//for reference only send locations > get wd > send back delta map
