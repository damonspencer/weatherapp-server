package main

import (
	"container/list"
	"log"
	"net/http"
	"weatherapp-server/weatherdata"

	"github.com/googollee/go-socket.io"
)

var clientmap map[string]string
var locationmap map[string]list.List

var locationslice []string

func contains(list []string, elem string) bool {
	for _, s := range list {
		if s == elem {
			return true
		}
	}
	return false
}

func getkeys(m map[string]list.List) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
func getindex(l list.List, query string) (match *list.Element) {
	match = l.Front()
	for i := l.Front(); i != nil && i.Value != query; i = i.Next() {
		match = i
	}
	if match.Value == query {
		return match
	}
	return nil
}
func socketserver(args []string) { //args currently not used here
	clientmap = make(map[string]string)
	//map of client to location of weather data requested,
	locationmap = make(map[string]list.List)
	//map of location of weather data requested to clients requesting it
	//really good for sending out weather data

	//all the wdc wants is the locations,
	//so lets process the map instad of sending the whole thing
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(socket socketio.Socket) {
		log.Println("Client", socket.Id(), "connected")
		socket.Join("users")

		socket.On("location update", func(newlocation string) {
			//make sure to provide the global location
			//of the city you want the weather of!
			//City:Threedigitcc
			log.Println("Client", socket.Id(), "updated location to", newlocation)
			socket.Join("wd")
			oldlocation := clientmap[socket.Id()] //clientmap is very useful for this
			clientmap[socket.Id()] = newlocation

			//remove from old map value
			oldlist := locationmap[oldlocation]
			elem := getindex(oldlist, socket.Id())
			if elem != nil {
				oldlist.Remove(elem)
				locationmap[oldlocation] = oldlist
			}

			newlist := locationmap[newlocation]
			newlist.PushBack(socket.Id())
			locationmap[newlocation] = newlist
			//add to new map value
			locationslice = getkeys(locationmap)
		})

		socket.On("newwd", func(wd map[string]weatherdata.Weatherdata) {
			for key, value := range wd {
				clientlist := locationmap[key]
				for e := clientlist.Front(); e != nil; e = e.Next() {
					s := e.Value.(string)
					socket.BroadcastTo(s, "newweatherdata", value)
				}
			}
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
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
