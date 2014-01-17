package main

import (
	"fmt"
	_ "log"
)

func main() {
	TestMe()
	fmt.Println("searching for devices to play to...")
	fmt.Println("now controlling [TV]UE04E20RRZZE")

	// commands := []string{"control", "play", "pause", "next", "prev", "search"}

	var d string
	for {
		_, err := fmt.Scanln(&d)
		if err != nil {
			break
		}
		fmt.Println(d)
	}

	// fmt.Println("Hello, world")

	// Parse("haggis")
	// CreateSoapRequest("a", "b", "c")
	// log.Fatal(fileserver.Serve(
	//  ":8080",
	//  "/Users/leetreveil/Downloads/Game.of.Thrones.S03E06.HDTV.x264-2HD.mp4",
	// ))

	// GetProtocolInfo()
}
