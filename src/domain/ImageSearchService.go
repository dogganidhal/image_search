package domain

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"log"
)

const (
	esHost              = "http://localhost:9200"
	imageLabelIndexName = "searchable_images_label"
	imageIndexName      = "searchable_images"
)

type ImageSearchService struct{}

func (E ImageSearchService) GetImagesByKeyword(keyword string, pagination Pagination) PaginatedResource {
	/* Search Query
	{
	  "query": {
	    "term": {
	      "label": $keyword
	    }
	  }
	}
	*/
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"label": keyword,
			},
		},
	}
	applyPagination(searchQuery, pagination)
	hits, total := executeElasticSearchQuery(imageLabelIndexName, searchQuery)
	return fetchImages(hits, pagination, total)
}

func (E ImageSearchService) GetImagesByRegex(regex string, pagination Pagination) PaginatedResource {
	/* Regex Query
	{
	  "query": {
	    "regexp": {
	      "label": {
	        "value": $regex,
	        "flags": "ALL",
	        "case_insensitive": true
	      }
	    }
	  }
	}
	*/
	regexQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"regexp": map[string]interface{}{
				"label": map[string]interface{}{
					"value":            regex,
					"flags":            "ALL",
					"case_insensitive": true,
				},
			},
		},
	}
	applyPagination(regexQuery, pagination)
	hits, total := executeElasticSearchQuery(imageLabelIndexName, regexQuery)
	return fetchImages(hits, pagination, total)
}

func (E ImageSearchService) GetImagesByKeywords(keywords []string, pagination Pagination) PaginatedResource {
	/* Search Query
	{
	  "searchQuery": {
	    "terms": {
	      "label": $keyword
	    }
	  }
	}
	*/
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"terms": map[string]interface{}{
				"label": keywords,
			},
		},
	}
	applyPagination(searchQuery, pagination)
	hits, total := executeElasticSearchQuery(imageLabelIndexName, searchQuery)
	return fetchImages(hits, pagination, total)
}

func executeElasticSearchQuery(index string, query map[string]interface{}) ([]interface{}, int64) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{esHost}})
	if err != nil {
		log.Fatalln(err)
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	response, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer response.Body.Close()
	if response.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				response.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	hits := result["hits"].(map[string]interface{})["hits"]
	total := result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"]
	return hits.([]interface{}), int64(total.(float64))
}

func fetchImages(hits []interface{}, pagination Pagination, total int64) PaginatedResource {
	imageIds := make([]string, 0)
	for _, hit := range hits {
		imageLabel := hit.(map[string]interface{})["_source"].(map[string]interface{})
		imageIds = append(imageIds, imageLabel["id"].(string))
	}
	/* Fetch Image Query
	{
	  "searchQuery": {
	    "terms": {
	      "id": $imageIds
	    }
	  }
	}
	*/
	fetchImagesQuery := map[string]interface{}{
		"size": pagination.Size,
		"query": map[string]interface{}{
			"terms": map[string]interface{}{
				"id": imageIds,
			},
		},
	}
	hits, _ = executeElasticSearchQuery(imageIndexName, fetchImagesQuery)
	return extractImageDocuments(hits, pagination, total)
}

func extractImageDocuments(hits []interface{}, pagination Pagination, total int64) PaginatedResource {
	items := make([]interface{}, 0)
	for _, hit := range hits {
		var image ImageDocument
		imageMap := hit.(map[string]interface{})["_source"]
		imageJson, err := json.Marshal(imageMap)
		if err != nil {
			log.Fatalln(err)
		}
		err = json.Unmarshal(imageJson, &image)
		if err != nil {
			log.Fatalln(err)
		}
		items = append(items, image)
	}
	return PaginatedResource{
		Total:      total,
		Pagination: pagination,
		Items:      items,
	}
}

func applyPagination(query map[string]interface{}, pagination Pagination) {
	query["from"] = pagination.From
	query["size"] = pagination.Size
}
