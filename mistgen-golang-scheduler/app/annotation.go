package app

import (
	"golang-mistgen/utils"
	"time"

	log "github.com/sirupsen/logrus"

	gapi "github.com/grafana/grafana-api-golang-client"
)

func grafanaAnnotate(config utils.GrafanaConfig, text string, timeStart time.Time, timeFinish time.Time) {
	log.Infof("sending annotation to grafana %v", text)
	clientConfig := gapi.Config{APIKey: config.ApiKey}
	client, err := gapi.New(config.Url, clientConfig)
	if err != nil {
		log.Error(err)
	}
	annotation := &gapi.Annotation{
		DashboardID: int64(config.DashboardId),
		Time:        makeEpoche(timeStart),
		TimeEnd:     makeEpoche(timeFinish),
		Text:        text,
		IsRegion:    true,
	}
	_, err = client.NewAnnotation(annotation)
	if err != nil {
		log.Error(err)
	}
	log.Info("annotation sent")
}

func makeEpoche(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
