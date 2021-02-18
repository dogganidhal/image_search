package web

import (
	"encoding/json"
	"ndogga.com/image_search/domain"
	"net/http"
	"strings"
)

type Controller struct{}

func (controller Controller) GetImagesByKeyword(writer http.ResponseWriter, request *http.Request) {
	imageSearchService := domain.ImageSearchService{}
	uriSlices := strings.Split(request.RequestURI, "/")
	keyword := uriSlices[len(uriSlices)-1]
	images := imageSearchService.GetImagesByKeyword(keyword)
	err := json.NewEncoder(writer).Encode(images)
	if err != nil {
		panic(err)
	}
}

func (controller Controller) GetImagesByRegex(writer http.ResponseWriter, request *http.Request) {
	imageSearchService := domain.ImageSearchService{}
	uriSlices := strings.Split(request.RequestURI, "/")
	regex := uriSlices[len(uriSlices)-1]
	images := imageSearchService.GetImagesByRegex(regex)
	err := json.NewEncoder(writer).Encode(images)
	if err != nil {
		panic(err)
	}
}

func (controller Controller) GetImagesByKeywords(writer http.ResponseWriter, request *http.Request) {
	imageSearchService := domain.ImageSearchService{}
	uriSlices := strings.Split(request.RequestURI, "/")
	keywords := uriSlices[len(uriSlices)-1]
	images := imageSearchService.GetImagesByKeywords(strings.Split(keywords, ","))
	err := json.NewEncoder(writer).Encode(images)
	if err != nil {
		panic(err)
	}
}
