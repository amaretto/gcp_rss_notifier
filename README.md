# gcp_rss_notifier
## description
gcp_rss_notifier get GCP RSS Feeds and notify user on slack

## Installation

- updatedb

Prepare database named "gcp-rss"


```shell
gcloud functions deploy UpdateDB --runtime go111 --trigger-topic updatedb --set-env-vars PROJECT_NAME=YOUR_PROJECT_NAME
```

- notifyinfo
```shell
gcloud functions deploy NotifyInfo --runtime go111 --trigger-topic notifyinfo --set-env-vars PROJECT_NAME=YOUR_PROJECT_NAME,SLACK_TOKEN=YOUR_SLACK_TOKEN,SLACK_USERNAME=YOUR_SLACK_USERNAME
