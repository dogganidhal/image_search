package main

import "ndogga.com/image_search/web"

func main() {
	webServer := web.Server{Port: 8080}
	webServer.Listen()
}
