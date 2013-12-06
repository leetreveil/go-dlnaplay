package main

import (
	_ "./fileserver"
	"./ssdp"
	_ "log"
)

func main() {
	ssdp.TestMe()
	// log.Fatal(fileserver.Serve(
	// 	":8080",
	// 	"/Users/leetreveil/Downloads/Game.of.Thrones.S03E06.HDTV.x264-2HD.mp4",
	// ))
}
