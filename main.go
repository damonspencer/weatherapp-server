package main

import (
	"log"
	"os"
	"weatherapp-server/tests"
)

func main() {
	args := os.Args[1:]
	//if contains(args, "test") { //testing mode only atm
	log.Println("Weather server starting up in test mode")
	//} else {
	log.Println("Weather server starting up with args:", args)
	//}
	b := make(chan int)
	go tests.Fakewdhttpserver(b) //oh well, now I have to build instead of go run...
	go wdclient(args, b)
	socketserver(args)
}
