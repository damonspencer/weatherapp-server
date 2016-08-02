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
	go tests.Fakewdhttpserver()
	go wdclient(args)
	socketserver(args)
}
