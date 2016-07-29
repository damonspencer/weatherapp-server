package main

import (
	"fmt"
	"time"
)

var weathermap map[string]string

func wdclient() {
	weathermap = make(map[string]string)
	//map of location to the actual weather at that place
	//make sure the socket server is running before taking its data
	for { //run until the program exits
		locationchannel := make(chan []string)
		locationchannel <- locationarray
		locations := <-locationchannel
		fmt.Println(locations)
		time.Sleep(1000 * time.Millisecond)
	}

}

//for reference only send locations > get wd > send back delta map
