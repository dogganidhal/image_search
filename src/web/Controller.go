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
	keyword := request.URL.Query().Get("keyword")
	images := imageSearchService.GetImagesByKeyword(keyword)
	err := json.NewEncoder(writer).Encode(images)
	if err != nil {
		panic(err)
	}
}

func (controller Controller) GetImagesByRegex(writer http.ResponseWriter, request *http.Request) {
	imageSearchService := domain.ImageSearchService{}
	regex := request.URL.Query().Get("regex")
	images := imageSearchService.GetImagesByRegex(regex)
	err := json.NewEncoder(writer).Encode(images)
	if err != nil {
		panic(err)
	}
}

func (controller Controller) GetImagesByKeywords(writer http.ResponseWriter, request *http.Request) {
	imageSearchService := domain.ImageSearchService{}
	keywords := request.URL.Query().Get("keywords")
	images := imageSearchService.GetImagesByKeywords(strings.Split(keywords, ","))
	err := json.NewEncoder(writer).Encode(images)
	if err != nil {
		panic(err)
	}
}
