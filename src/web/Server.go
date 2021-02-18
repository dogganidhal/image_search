package web

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Port int
}

func (server Server) Listen() {
	addr := fmt.Sprintf(":%d", server.Port)
	controller := Controller{}

	http.HandleFunc("/search/keyword/", controller.GetImagesByKeyword)
	http.HandleFunc("/search/regex/", controller.GetImagesByRegex)
	http.HandleFunc("/search/keywords/", controller.GetImagesByKeywords)

	log.Fatal(http.ListenAndServe(addr, nil))
}
