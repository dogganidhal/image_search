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

func (E ImageSearchService) GetImagesByKeyword(keyword string) []ImageDocument {
	/* Search Query
	{
	  "searchQuery": {
	    "term": {
	      "label": {
	        "value": $keyword
	      }
	    }
	  }
	}
	*/
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"label": map[string]interface{}{
					"value": keyword,
				},
			},
		},
	}
	hits := executeElasticSearchQuery(imageLabelIndexName, searchQuery)
	var imageIds []string
	for _, hit := range hits {
		imageLabel := hit.(map[string]interface{})["_source"].(map[string]interface{})
		imageIds = append(imageIds, imageLabel["id"].(string))
	}
	/* Fetch Image Query
	{
	  "searchQuery": {
	    "term": {
	      "id": {
	        "value": $imageIds
	      }
	    }
	  }
	}
	*/
	fetchImagesQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"terms": map[string]interface{}{
				"id": imageIds,
			},
		},
	}
	hits = executeElasticSearchQuery(imageIndexName, fetchImagesQuery)
	var result []ImageDocument
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
		result = append(result, image)
	}
	return result
}

func (E ImageSearchService) GetImagesByRegex(regex string) []ImageDocument {
	panic("implement me")
}

func (E ImageSearchService) GetImagesByKeywords(keywords []string) []ImageDocument {
	panic("implement me")
}

func executeElasticSearchQuery(index string, query map[string]interface{}) []interface{} {
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
	return hits.([]interface{})
}
