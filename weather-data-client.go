package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"weatherapp-server/weatherdata"

	"github.com/zhouhui8915/go-socket.io-client"
)

const fakewdurl = "http://localhost:5555/fakewd"
const socketserverurl = "http://localhost:5000"

var oldweathermap map[string]weatherdata.Weatherdata
var newweathermap map[string]weatherdata.Weatherdata
var deltaweathermap map[string]weatherdata.Weatherdata
var oldwd []weatherdata.Weatherdata

func wdclient(args []string) {
	test := false
	if contains(args, "test") {
		test = true
	}
	log.Println(test)
	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	time.Sleep(time.Second * 10)
	//wait for the fake weather data server to start up
	//there's got to be a better way....
	//TODO find a better way

	oldweathermap = make(map[string]weatherdata.Weatherdata)
	newweathermap = make(map[string]weatherdata.Weatherdata)
	//map of locations to the actual weather at that place
	//we need to have a map to get the delta of later
	newweathermap = testwdsetup(fakewdurl)
	client, err := socketio_client.NewClient(socketserverurl, opts)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("The wdc is connected! The client id above is probably the wdc")
	}
	//emit the initial weather data
	client.Emit("newwd", newweathermap)
	locationchannel := make(chan []string)
	for { //run until the program exits
		locationchannel <- locationslice
		locations := <-locationchannel
		//only test mode is implemented
		newweathermap = gettestwd(fakewdurl, locations)
		deltaweathermap = getdeltamap(oldweathermap, newweathermap)
		oldweathermap = newweathermap
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
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(body)
	//fmt.Println(string(body))
	var newwdslice []weatherdata.Weatherdata
	err = json.Unmarshal(body[:len(body)-23], &newwdslice)
	//"%!(EXTRA string=fakewd)" appears at the end, so discard it
	//a bad way of extracting a json
	//from a byte array of a json turned into a
	//string turned back into a byte array...
	//maybe I could just discard everything after "%"?
	//or maybe I could make a better http server
	//so I wouldn't have to do this in the first place...
	//fmt.Println(newwdslice)
	newweathermap := make(map[string]weatherdata.Weatherdata)
	for i := range newwdslice {
		globallocation := newwdslice[i].City + ":" + newwdslice[i].Threedigitcc
		if contains(locations, globallocation) {
			newweathermap[globallocation] = newwdslice[i]
		}
	}
	return newweathermap
}

func testwdsetup(url string) map[string]weatherdata.Weatherdata {
	//the same as gettestwd minus the locations
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var newwdslice []weatherdata.Weatherdata
	err = json.Unmarshal(body[:len(body)-23], &newwdslice)
	newweathermap := make(map[string]weatherdata.Weatherdata)
	for i := range newwdslice {
		globallocation := newwdslice[i].City + ":" + newwdslice[i].Threedigitcc
		newweathermap[globallocation] = newwdslice[i]
	}
	return newweathermap
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

	return deltamap //I think this is right
}
