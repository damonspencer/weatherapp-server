package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/googollee/go-socket.io"
)

var clientmap map[string]string
var locationmap map[string]string

//var locationlist list
var locationarray []string

//var locationlist list

//I don't think this needs to be global
func contains(list []string, elem string) bool {
	for _, s := range list {
		if s == elem {
			return true
		}
	}
	return false
}

func getkeys(m map[string]string) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

//func handleerror(){}
func socketserver(args []string) { //args currently not used here
	clientmap = make(map[string]string)
	//map of client to location of weather data requested,
	//I'm not sure if I'll use this, the reverse map seems more useful
	locationmap = make(map[string]string)
	//map of location of weather data requested to clients
	//	var locationlist map[string]string = make(map[string]string)
	//all the wdc wants is the locations,
	//so lets process the map instad of sending the whole thing
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(socket socketio.Socket, location string) {
		log.Println("Client", socket.Id(),
			"connected wanting the weather of", location)
		clientmap[socket.Id()] = location
		locationmap[location] += socket.Id() + ":"
		//I think socket ids are hashes, right, so ":" should be a good seperator
		socket.Join("users") //for general messages to all users
		//TODO figure out how to get user count
		//log.Println()
		socket.On("location update", func(loc string) {
			log.Println("Client", socket.Id(), "updated location to", loc)
			oldlocation := clientmap[socket.Id()]
			clientmap[socket.Id()] = loc

			//when we use the reverse map instead of the regular clientmap,
			// we are betting clients will update their location less
			//than we poll the server, because updating their location is
			//computationally expensive,
			//but only happens when they update their location.
			// while getting a location list from the regular map is computationally
			//expensive but happens every time we need to poll
			//the server if we use the regular map. I think using the
			//reverse map is a safe bet at minute polling intervals.
			//this is only correct with a relatively small number of clients.
			//It would be the opposite for larger numbers of clients.
			//Disclaimer:
			//This is all completely theoretical and came out my head at 4am

			strings.Replace(locationmap[oldlocation], socket.Id()+":", "", -1)
			//remove from old value
			locationmap[loc] += socket.Id() + ":"
			//add to new value
			locationarray = getkeys(locationmap)
		})

		socket.On("disconnection", func() {
			log.Println("Client", socket.Id(), "disconnected")
			delete(clientmap, socket.Id())
		})

	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
	http.Handle("/socket.io/", server)
	http.Handle("/socketserver", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5555...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
