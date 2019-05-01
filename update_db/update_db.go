package main

import (
	"context"
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
)

// FeedURL is GCP RSS Feed URL
const FeedURL string = "https://status.cloud.google.com/feed.atom"

// GcpRss divide RSS Feed into some parts
type GcpRss struct {
	EntryTitle []string `xml:"entry>title"`
	Updated    []string `xml:"entry>updated"`
	Content    []string `xml:"entry>content"`
}

func main() {
	var prjName string
	flag.StringVar(&prjName, "prjName", "", "GCP project name")
	flag.Parse()

	//Prepare client
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, prjName)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	gr, err := getGcpRss(FeedURL)
	if err != nil {
		log.Fatalf("Log: %v", err)
		return
	}

	layout := "2006-01-02T15:04:05Z"

	feedCount := len(gr.EntryTitle)
	for count := 0; count < feedCount; count++ {
		// Convert updated to unix time as ID
		t, err := time.Parse(layout, gr.Updated[count])
		if err != nil {
			log.Fatalf("Log: %v", err)
		}

		//Check the record already exist or not
		if isExistRecord(ctx, t.Unix(), client) {
			continue
		}

		// Parse elements
		re := regexp.MustCompile(`([A-Z]*):\s([A-Za-z]*)\s([0-9]*)\s-\s(.*)`)
		result := re.FindStringSubmatch(gr.EntryTitle[count])
		// result:1->STATUS,3->INCIDENT_NO,4->TITLE

		registerRecord(ctx, t.Unix(), result[1], gr.Updated[count], result[3], result[4], gr.Content[count], client)
	}
}

// Get GCP RSS
func getGcpRss(feed string) (p *GcpRss, err error) {
	res, err := http.Get(feed)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	gr := new(GcpRss)
	err = xml.Unmarshal(b, &gr)

	return gr, err
}

// check Record is already exists or not
func isExistRecord(ctx context.Context, id int64, c *firestore.Client) bool {
	doc := c.Collection("gcp-rss").Doc(strconv.FormatInt(id, 10))
	_, err := doc.Get(ctx)
	if err != nil {
		//log.Fatalf("Failed to open doc : %v", err)
		//ToDo:Check type of error
		return false
	}
	return true
}

//
func registerRecord(ctx context.Context, id int64, status, updated, incidentNo, title, detail string, c *firestore.Client) {
	_, err := c.Collection("gcp-rss").Doc(strconv.FormatInt(id, 10)).Set(ctx, map[string]interface{}{
		"STATUS":      status,
		"UPDATED":     updated,
		"INCIDENT_NO": incidentNo,
		"TITLE":       title,
		"DETAIL":      detail,
	})
	if err != nil {
		log.Fatalf("Log: %v", err)
	}
}
