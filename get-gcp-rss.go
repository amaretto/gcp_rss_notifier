package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var FEED_URL string = "https://status.cloud.google.com/feed.atom"

type GcpRss struct {
	Title      string   `xml:"title"`
	EntryTitle []string `xml:"entry>title"`
	Updated    []string `xml:"entry>updated"`
	Content    []string `xml:"entry>content"`
}

func main() {
	gr, err := getGcpRss(FEED_URL)
	if err != nil {
		log.Fatalf("Log: %v", err)
		return
	}

	fmt.Println(gr.Title)
	feedCount := len(gr.EntryTitle)
	fmt.Println("===========================")
	for count := 0; count < feedCount; count++ {
		fmt.Printf("[Title]\n%s\n\n", gr.EntryTitle[count])
		fmt.Printf("[Updated]\n%s\n\n", gr.Updated[count])
		fmt.Printf("[Detail]\n%s\n\n", gr.Content[count])
		fmt.Println("--------------------------")
	}
}

// Get Weather Hacaks
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
