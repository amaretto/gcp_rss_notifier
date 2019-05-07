package p

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/nlopes/slack"
	"google.golang.org/api/iterator"
)

// ToDo: Need to refactor // - unite format of log  and error
// - rename variable

// GcpRssInfo have information of RSS from GCP
type GcpRssInfo struct {
	Status     string
	Updated    string
	IncidentNo string
	Title      string
	Detail     string
}

// PubSubMessage accept message from Cloud Pub/Sub
type PubSubMessage struct {
	Data []byte `json:data`
}

// NotifyInfo notify new GCP RSS(updated after last notification) for users
func NotifyInfo(ctx context.Context, m PubSubMessage) error {

	//ctx := context.Background()
	r, err := getNewInfo(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if len(r) == 0 {
		return nil
	}

	var lastUpd string
	var messageString string

	messageString = "===========================\n"
	for _, gri := range r {
		messageString += "[Title]\n" + gri.Title + "(" + gri.IncidentNo + ")\n\n"
		messageString += "[Status]\n" + gri.Status + "\n\n"
		messageString += "[Update]\n" + gri.Updated + "\n\n"
		messageString += "[Detail]\n" + gri.Detail + "\n\n"
		messageString += "--------------------------\n"

		lastUpd = gri.Updated
	}
	SendMessage(messageString)

	// ToDo : update update-time collection
	updateLastUpdTime(ctx, lastUpd)
	return nil
}

// get new GCP RSS Info from firestore
func getNewInfo(ctx context.Context) (p []*GcpRssInfo, err error) {
	prjName := os.Getenv("PROJECT_NAME")
	client, err := firestore.NewClient(ctx, prjName)
	if err != nil {
		log.Fatalf("Failed to create firestore client: %v", err)
	}
	defer client.Close()

	// get last updated timestamp
	doc := client.Collection("update-time").Doc("last-updated")
	docsnap, err := doc.Get(ctx)
	if err != nil {
		log.Fatalf("Faild to open doc : %v", err)
	}
	lastUpd := docsnap.Data()["UPDATED"]

	// get new info from firestore
	// ToDo :use environment variable for collection name
	iter := client.Collection("gcp-rss").Where("UPDATED", ">", lastUpd).Documents(ctx)

	for {
		docsnap, err = iter.Next()
		if err == iterator.Done {
			err = nil
			break
		}
		if err != nil {
			log.Fatalf("Faild to iterate: %v", err)
		}

		ni := new(GcpRssInfo)
		ni.Status = docsnap.Data()["STATUS"].(string)
		ni.Updated = docsnap.Data()["UPDATED"].(string)
		ni.IncidentNo = docsnap.Data()["INCIDENT_NO"].(string)
		ni.Title = docsnap.Data()["TITLE"].(string)
		ni.Detail = docsnap.Data()["DETAIL"].(string)
		p = append(p, ni)
	}

	return p, err
}

func updateLastUpdTime(ctx context.Context, lastUpd string) (err error) {
	// ToDo : Avoid writing same code in getNewInfo
	prjName := os.Getenv("PROJECT_NAME")
	client, err := firestore.NewClient(ctx, prjName)
	if err != nil {
		log.Fatalf("Failed to create firestore client: %v", err)
	}
	defer client.Close()

	_, err = client.Collection("update-time").Doc("last-updated").Set(ctx, map[string]interface{}{
		"UPDATED": lastUpd,
	})
	if err != nil {
		return err
	}
	return nil
}

// SendMessage send message to Slack
func SendMessage(message string) {
	// ToDo: return err
	// ToDo: this method use old library. it is needed adapting latest version.
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	username := os.Getenv("SLACK_USERNAME")

	message = "<" + username + "> \n" + message
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Text: message,
	}
	params.Attachments = []slack.Attachment{attachment}
	params.AsUser = true

	api.PostMessage("#gcp_rss", "", params)
}
