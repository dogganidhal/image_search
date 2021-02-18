package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"github.com/cenkalti/backoff/v4"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/google/uuid"
	"io"
	"log"
	"ndogga.com/image_search/domain"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	esHost              = "http://localhost:9200"
	imageLabelIndexName = "searchable_images_label"
	imageIndexName      = "searchable_images"
)

type imageLabelName struct {
	ImageId   string
	LabelName string
}

func CreateElasticSearch() *elasticsearch.Client {
	retryBackoff := backoff.NewExponentialBackOff()
	esConfig := elasticsearch.Config{
		Addresses:     []string{esHost},
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		MaxRetries: 5,
	}
	es, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		panic(err)
	}
	return es
}

func readImageLabels(confidenceThreshold float64) []imageLabelName {
	csvFile, err := os.Open("./data/image-labels.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	csvReader := csv.NewReader(csvFile)
	var imageLabels []imageLabelName
	for {
		lineRecords, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Couldn't read the csv file", err)
		}
		confidence, err := strconv.ParseFloat(lineRecords[3], 64)
		if confidence >= confidenceThreshold {
			imageLabel := imageLabelName{
				ImageId:   lineRecords[0],
				LabelName: lineRecords[2],
			}
			imageLabels = append(imageLabels, imageLabel)
		}
	}
	return imageLabels
}

func readClassDescriptions() map[string]string {
	classDescriptions := make(map[string]string)
	csvFile, err := os.Open("./data/class-descriptions.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	csvReader := csv.NewReader(csvFile)
	for {
		lineRecords, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Couldn't read the csv file", err)
		}
		classDescriptions[lineRecords[0]] = lineRecords[1]
	}
	return classDescriptions
}

func indexImageLabels(es *elasticsearch.Client, imageLabels []imageLabelName, labelDescriptions map[string]string) {
	bulkConfig := esutil.BulkIndexerConfig{
		Client: es,
		Index:  imageLabelIndexName,
	}
	bulkIndexer, err := esutil.NewBulkIndexer(bulkConfig)
	if err != nil {
		log.Fatalln(err)
	}
	start := time.Now().UTC()
	for _, imageLabel := range imageLabels {
		document := domain.ImageLabelDocument{
			Id:    imageLabel.ImageId,
			Label: labelDescriptions[imageLabel.LabelName],
		}
		documentJsonBytes, err := json.Marshal(document)
		if err != nil {
			log.Fatalln(err)
		}
		documentId := uuid.New().String()
		item := esutil.BulkIndexerItem{
			Action:     "index",
			Index:      imageLabelIndexName,
			DocumentID: documentId,
			Body:       strings.NewReader(string(documentJsonBytes)),
		}
		err = bulkIndexer.Add(context.Background(), item)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if err := bulkIndexer.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
	stats := bulkIndexer.Stats()
	duration := time.Since(start)

	if stats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(stats.NumFlushed)),
			humanize.Comma(int64(stats.NumFailed)),
			duration.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(duration/time.Millisecond)*float64(stats.NumFlushed))),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(stats.NumFlushed)),
			duration.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(duration/time.Millisecond)*float64(stats.NumFlushed))),
		)
	}
}

func indexImages(es *elasticsearch.Client) {
	csvFile, err := os.Open("./data/image-descriptions.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	csvReader := csv.NewReader(csvFile)
	bulkConfig := esutil.BulkIndexerConfig{
		Client: es,
		Index:  imageLabelIndexName,
	}
	bulkIndexer, err := esutil.NewBulkIndexer(bulkConfig)
	if err != nil {
		log.Fatalln(err)
	}
	start := time.Now().UTC()
	for {
		lineRecords, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Couldn't read the csv file", err)
		}
		originalSize, err := strconv.ParseFloat(lineRecords[8], 64)
		if err != nil {
			log.Fatalln(err)
		}
		//ImageID,Subset,OriginalURL,OriginalLandingURL,License,AuthorProfileURL,Author,Title,OriginalSize,OriginalMD5,Thumbnail300KURL,Rotation
		//0		 ,1		,2			,3				   ,4	   ,5				,6	   ,7	 ,8			  ,9		  ,10			   ,11
		image := domain.ImageDocument{
			Id:                 lineRecords[0],
			OriginalUrl:        lineRecords[2],
			OriginalLandingURL: lineRecords[3],
			AuthorProfileURL:   lineRecords[5],
			Author:             lineRecords[6],
			Title:              lineRecords[7],
			OriginalSize:       originalSize,
			OriginalMD5:        lineRecords[9],
			Thumbnail300KUrl:   lineRecords[10],
		}
		imageJsonBytes, err := json.Marshal(image)
		if err != nil {
			log.Fatalln(err)
		}
		documentId := uuid.New().String()
		item := esutil.BulkIndexerItem{
			Action:     "index",
			Index:      imageIndexName,
			DocumentID: documentId,
			Body:       strings.NewReader(string(imageJsonBytes)),
		}
		err = bulkIndexer.Add(context.Background(), item)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if err := bulkIndexer.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
	stats := bulkIndexer.Stats()
	duration := time.Since(start)

	if stats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(stats.NumFlushed)),
			humanize.Comma(int64(stats.NumFailed)),
			duration.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(duration/time.Millisecond)*float64(stats.NumFlushed))),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(stats.NumFlushed)),
			duration.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(duration/time.Millisecond)*float64(stats.NumFlushed))),
		)
	}
}

func main() {
	es := CreateElasticSearch()
	log.Println(es.Info())
	log.Printf("Successfully created Elastic Search client")
	imageLabels := readImageLabels(0.8)
	classDescriptions := readClassDescriptions()
	indexImageLabels(es, imageLabels, classDescriptions)
	indexImages(es)
	log.Printf("Found %d image labels", len(imageLabels))
}
