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
	server, err := socketio.NewServer(nil)

	//TODO keep track of users
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Weather server starting up")

	server.On("connection", func(socket socketio.Socket) {

		log.Println("User connected") //maybe add user name?
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

		socket.On("disconnection", func(clientid string) {
			log.Println(clientid, "disconnected") //maybe add client id?
		})

	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
