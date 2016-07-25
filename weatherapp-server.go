package main

import (
	"log"
	"net/http"
	"os"

	"github.com/googollee/go-socket.io"
)

func contains(list []string, elem string) bool {
	for _, s := range list {
		if s == elem {
			return true
		}
	}
	return false
}
func main() {
	args := os.Args[1:]
	if contains(args, "test") {
		log.Println("Weather server starting up in test mode")
	} else if args != nil {
		log.Println("Weather server starting up with args", args)
	} else {
		log.Println("Weather server starting up")
	}
	var Clientmap map[string]string = make(map[string]string)
	//map of client to location of weather data requested
	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(socket socketio.Socket, location string) {
		log.Println("Client", socket.Id(),
			"connected wanting the weather of", location)
		Clientmap[socket.Id()] = location
		socket.Join("users") //for general messages to all users
		//log.Println("there are", get_user_count(), get_user_list())
		log.Println()
		//socket.Join("locationupdate")
		//socket.Join("sendweatherdata")
		//socket.On("locationupdate message", func(msg string) {
		//log.Println("emiting:", socket.Emit("locationupdate message", msg))
		//log.Println("Received location update from", clientid)
		//socket.Emit("locationupdate message", msg))
		//socket.Emit("locationupdate", "locationupdate message", msg)
		//aknowlege we recived their location update request
		//})

		//socket.On("locationupdate message", func(msg string) {
		//log.Println("emiting:", socket.Emit("sendweatherdata message", msg))
		//socket.Emit("weatherdata", "weatherdata message", msg)
		//update the client with new weather data
		//})
		socket.On("location update", func(location string) {
			log.Println("Client", socket.Id(), "updated location to", location)
			Clientmap[socket.Id()] = location
		})

		socket.On("disconnection", func() { //(socket socketio.Socket) {
			log.Println("Client", socket.Id(), "disconnected")
			delete(Clientmap, socket.Id())
		})

	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
	//setInterval(func() {
	//var randomClient
	//if (clients.length > 0) {
	//randomClient = Math.floor(Math.random() * clients.length)
	//clients[randomClient].emit('foo', sequence++)
	//}
	//}, 1000)
	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
