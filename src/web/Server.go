package web

import (
	"fmt"
	"github.com/rs/cors"
	"log"
	"net/http"
)

type Server struct {
	Port int
}

func (server Server) Listen() {
	addr := fmt.Sprintf(":%d", server.Port)
	controller := Controller{}

	mux := http.NewServeMux()

	mux.HandleFunc("/search/keyword/", controller.GetImagesByKeyword)
	mux.HandleFunc("/search/regex/", controller.GetImagesByRegex)
	mux.HandleFunc("/search/keywords/", controller.GetImagesByKeywords)

	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(addr, handler))
}
