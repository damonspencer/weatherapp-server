package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if contains(args, "test") {
		log.Println("Weather server starting up in test mode")
	} else if args != nil {
		log.Println("Weather server starting up with args", args)
	} else {
		log.Println("Weather server starting up")
	}
	go wdclient(args)
	socketserver(args)
}
