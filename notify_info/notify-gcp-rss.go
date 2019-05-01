package main

import (
	"context"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
)

var FEED_URL string = "https://status.cloud.google.com/feed.atom"

type GcpRss struct {
	Title      string   `xml:"title"`
	EntryTitle []string `xml:"entry>title"`
	Updated    []string `xml:"entry>updated"`
	Content    []string `xml:"entry>content"`
}

func main() {
	lambda.Start(HandlerRequest)
}

func HandlerRequest(ctx context.Context, params interface{}) (interface{}, error) {

	var messageString string

	gr, err := getGcpRss(FEED_URL)
	if err != nil {
		log.Fatalf("Log: %v", err)
		os.Exit(1)
	}
	messageString = gr.Title + "\n===========================\n"
	feedCount := len(gr.EntryTitle)
	for count := 0; count < feedCount; count++ {
		messageString += "[Title]\n" + gr.EntryTitle[count] + "\n\n"
		messageString += "[Update]\n" + gr.Updated[count] + "\n\n"
		messageString += "[Detail]\n" + gr.Content[count] + "\n\n"
		messageString += "--------------------------"
	}
	SendMessage(messageString)
	return params, nil
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

func SendMessage(message string) {
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	username := os.Getenv("SLACK_USERNAME")

	message = "<" + username + "> " + message
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Text: message,
	}
	params.Attachments = []slack.Attachment{attachment}
	params.AsUser = true

	api.PostMessage("#bot_project", "", params)
}
