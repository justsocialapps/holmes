package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/justsocialapps/holmes/analytics"
	"github.com/justsocialapps/holmes/assets"
	"github.com/justsocialapps/holmes/publisher"
	"github.com/justsocialapps/holmes/tracker"
	"gopkg.in/Shopify/sarama.v1"
)

//go:generate scripts/prepare_assets.sh
//go:generate go run scripts/include_assets.go
//go:generate gofmt -w assets/assets.go

func provideTrackingParams(params tracker.TrackingParams, handler func(params tracker.TrackingParams, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(params, w, r)
	}
}

func startServer(host string, port string) {
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Holmes %s running on %s:%s", version, host, port)
	log.Println(assets.Bannertxt)

	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func initLogging(logfileName string) {
	logFile := os.Stdout
	if logfileName != "" {
		var err error
		logFile, err = os.OpenFile(logfileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.SetOutput(logFile)
	sarama.Logger = log.New(logFile, "[Sarama] ", log.LstdFlags)
}

func main() {
	if printVersion {
		fmt.Println(version)
		return
	}

	initLogging(logfileName)

	trackingChannel := make(chan *tracker.TrackingObject, 10)
	trackingParams := tracker.TrackingParams{
		TrackingChannel: trackingChannel,
		AnonymizeIP:     anomyzeIP,
	}

	go publisher.Publish(trackingChannel, kafkaHost, "tracking")

	http.HandleFunc("/track", provideTrackingParams(trackingParams, tracker.Track))
	http.HandleFunc("/analytics.js", analytics.Analytics)

	startServer(listenHost, listenPort)
}
