package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/justsocialapps/holmes/assets"
	"github.com/justsocialapps/holmes/handlers"
	"github.com/justsocialapps/holmes/models"
	"github.com/justsocialapps/holmes/publisher"
)

//go:generate scripts/prepare_assets.sh
//go:generate go run scripts/include_assets.go

func provideTrackingChannel(trackingChannel chan<- *models.TrackingObject, handler func(trackingChannel chan<- *models.TrackingObject, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
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
	flag.Parse()

	baseUrl := *protocol + "://" + *host
	if !(*protocol == "https" && *proxyPort == "443") && !(*protocol == "http" && *proxyPort == "80") {
		baseUrl = baseUrl + ":" + *proxyPort
	}
	baseUrl = baseUrl + *proxyPath

	trackingChannel := make(chan *models.TrackingObject, 10)
	http.HandleFunc("/track", provideTrackingChannel(trackingChannel, handlers.Track))
	http.HandleFunc("/analytics.js", handlers.Analytics(baseUrl))
	go publisher.Publish(trackingChannel, kafkaHost, "tracking")

	startServer("localhost", listenPort)
}
