package web

import (
	"encoding/json"
	"log"
	"ndogga.com/image_search/domain"
	"net/http"
	"strconv"
	"strings"
)

type Controller struct{}

func (controller Controller) GetImagesByKeyword(writer http.ResponseWriter, request *http.Request) {
	imageSearchService := domain.ImageSearchService{}
	uriSlices := strings.Split(request.RequestURI, "/")
	keyword := extractQuery(uriSlices)
	pagination := extractPagination(request)
	images := imageSearchService.GetImagesByKeyword(keyword, pagination)
	err := json.NewEncoder(writer).Encode(images)
	if err != nil {
		panic(err)
	}
	enableCors(&writer)
}

func (controller Controller) GetImagesByRegex(writer http.ResponseWriter, request *http.Request) {
	imageSearchService := domain.ImageSearchService{}
	uriSlices := strings.Split(request.RequestURI, "/")
	regex := extractQuery(uriSlices)
	pagination := extractPagination(request)
	images := imageSearchService.GetImagesByRegex(regex, pagination)
	err := json.NewEncoder(writer).Encode(images)
	if err != nil {
		panic(err)
	}
	enableCors(&writer)
}

func (controller Controller) GetImagesByKeywords(writer http.ResponseWriter, request *http.Request) {
	imageSearchService := domain.ImageSearchService{}
	uriSlices := strings.Split(request.RequestURI, "/")
	keywords := extractQuery(uriSlices)
	pagination := extractPagination(request)
	images := imageSearchService.GetImagesByKeywords(strings.Split(keywords, ","), pagination)
	err := json.NewEncoder(writer).Encode(images)
	if err != nil {
		panic(err)
	}
	enableCors(&writer)
}

func extractPagination(request *http.Request) domain.Pagination {
	fromParameter := request.URL.Query().Get("from")
	sizeParameter := request.URL.Query().Get("size")
	var from, size int64
	var err error
	if len(fromParameter) > 0 {
		from, err = strconv.ParseInt(fromParameter, 10, 64)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		from = domain.DefaultPaginationFrom()
	}
	if len(sizeParameter) > 0 {
		size, err = strconv.ParseInt(sizeParameter, 10, 64)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		size = domain.DefaultPaginationSize()
	}
	return domain.Pagination{
		From: from,
		Size: size,
	}
}

func extractQuery(slices []string) string {
	lastUriPath := slices[len(slices)-1]
	lastUriPathSlices := strings.Split(lastUriPath, "?")
	return lastUriPathSlices[0]
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
