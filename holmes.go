package main

import (
	"flag"
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

const version string = "1.7.0-dev"

//go:generate scripts/prepare_assets.sh
//go:generate go run scripts/include_assets.go

func provideTrackingChannel(trackingChannel chan<- *tracker.TrackingObject, handler func(trackingChannel chan<- *tracker.TrackingObject, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(trackingChannel, w, r)
	}
}

func startServer(host string, port *string) {
	listener, err := net.Listen("tcp", host+":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Holmes running on " + host + ":" + *port)
	log.Println(assets.Bannertxt)

	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var protocol = flag.String("protocol", "https", "The protocol used to serve Holmes resources ('http' or 'https')")
	var host = flag.String("host", "localhost", "The host name used to reach Holmes")
	var proxyPort = flag.String("proxyPort", "3001", "The TCP port for reaching Holmes if Holmes is operated behind a reverse proxy.")
	var proxyPath = flag.String("proxyPath", "", "The base path for reaching Holmes if Holmes is operated behind a reverse proxy.")
	var listenPort = flag.String("listenPort", "3001", "The TCP port that Holmes listens on")
	var kafkaHost = flag.String("kafkaHost", "localhost:9092", "The Kafka host to consume messages from")
	var logfileName = flag.String("logfile", "", "The file to log messages to")
	var printVersion = flag.Bool("version", false, "Print Holmes version and exit")
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		return
	}

	var logFile *os.File
	if *logfileName == "" {
		logFile = os.Stdout
	} else {
		var err error
		logFile, err = os.OpenFile(*logfileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.SetOutput(logFile)
	sarama.Logger = log.New(logFile, "[Sarama] ", log.LstdFlags)

	baseURL := *protocol + "://" + *host
	if !(*protocol == "https" && *proxyPort == "443") && !(*protocol == "http" && *proxyPort == "80") {
		baseURL = baseURL + ":" + *proxyPort
	}
	baseURL = baseURL + *proxyPath

	trackingChannel := make(chan *tracker.TrackingObject, 10)
	http.HandleFunc("/track", provideTrackingChannel(trackingChannel, tracker.Track))
	http.HandleFunc("/analytics.js", analytics.Analytics(baseURL))
	go publisher.Publish(trackingChannel, kafkaHost, "tracking")

	startServer("localhost", listenPort)
}
