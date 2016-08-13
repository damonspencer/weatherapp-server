package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"weatherapp-server/weatherdata"

	"github.com/zhouhui8915/go-socket.io-client"
)

const fakewdurl = "http://localhost:5555/fakewd"
const socketserverurl = "http://localhost:5000"

func wdclient(args []string, b chan int) {
	//test := false
	//if contains(args, "test") {
	//	test = true
	//}
	//log.Println(test)

	//only test mode is implemented

	oldweathermap := make(map[string]weatherdata.Weatherdata)
	newweathermap := make(map[string]weatherdata.Weatherdata)
	deltaweathermap := make(map[string]weatherdata.Weatherdata)
	//map of locations to the actual weather at that place
	//we need to have maps to get the delta of later

	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	client, err := socketio_client.NewClient(socketserverurl, opts)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("The wdc is connected! The client id above is probably the wdc")
	}
	locationchannel := make(chan []string)
	locationchannel <- locationslice
	locations := <-locationchannel
	//wait until a client requests the weather before polling

	//wait for the fake weather data server to start up
	<-b
	newweathermap = gettestwd(fakewdurl, locations)
	//we'll wait until more than one request for the weather has been made before getting the delta

	//emit the initial weather data
	client.Emit("newwd", newweathermap)
	for { //run until the program exits
		locationchannel <- locationslice
		locations := <-locationchannel
		//only test mode is implemented
		newweathermap = gettestwd(fakewdurl, locations)
		deltaweathermap = getdeltamap(oldweathermap, newweathermap)
		oldweathermap = newweathermap
		//set the old one to the new one

		//emit deltaweathermap...
		client.Emit("newwd", deltaweathermap)
		time.Sleep(time.Second) //sleep for 1 second
	}

}

func gettestwd(url string, locations []string) map[string]weatherdata.Weatherdata {
	////func gettestwd(url string) map[string]weatherdata.Weatherdata {
	//returns a map of wd to locations

	//we're not using a channel because
	//we're trying to pretend this is an unrelated website
	resp, httperr := http.Get(url)
	handleerror(httperr)
	defer resp.Body.Close()
	body, ioerr := ioutil.ReadAll(resp.Body)
	handleerror(ioerr)
	//fmt.Println(body)
	//fmt.Println(string(body))
	var wdslice []weatherdata.Weatherdata
	jsonerr := json.Unmarshal(body[:len(body)-23], &wdslice)
	handleerror(jsonerr)
	//"%!(EXTRA string=fakewd)" appears at the end, so discard it
	//a bad way of extracting a json
	//from a byte array of a json turned into a
	//string turned back into a byte array...
	//maybe I could just discard everything after "%"?
	//or maybe I could make a better http server
	//so I wouldn't have to do this in the first place...
	//fmt.Println(newwdslice)
	newmap := make(map[string]weatherdata.Weatherdata)
	for i := range wdslice {
		globallocation := wdslice[i].City + ":" + wdslice[i].Threedigitcc
		if contains(locations, globallocation) {
			newmap[globallocation] = wdslice[i]
		}
	}
	return newmap
}

func getdeltamap(
	oldmap map[string]weatherdata.Weatherdata,
	newmap map[string]weatherdata.Weatherdata) map[string]weatherdata.Weatherdata {
	//four things could happen: we have more locations in the new map
	//(a client changed their location to somewhere new),
	//less locations in the new map, some key's value has changed
	//or a key's value hasn't changed
	deltamap := make(map[string]weatherdata.Weatherdata)
	for key, value := range newmap {
		if oldmap[key] != newmap[key] {
			//we don't have to test any keys that are not in the newmap,
			//because the clients needing the weatherdata
			//for that location have disconnected
			//new locations should be put in the deltamap
			//stuff that has changed should be in there too
			deltamap[key] = value
		}
	}

	return deltamap
}
func handleerror(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
