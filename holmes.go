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
func initLogging(logfileName *string) {
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
}

func main() {
	var protocol = flag.String("protocol", "https", "The protocol used to serve Holmes resources ('http' or 'https')")
	var host = flag.String("host", "localhost", "The host name used to reach Holmes")
	var proxyPort = flag.String("proxyPort", "3001", "The TCP port for reaching Holmes if Holmes is operated behind a reverse proxy.")
	var proxyPath = flag.String("proxyPath", "", "The base path for reaching Holmes if Holmes is operated behind a reverse proxy.")
	var listenHost = flag.String("listenHost", "", "The host Holmes listens on.")
	var listenPort = flag.String("listenPort", "3001", "The TCP port that Holmes listens on")
	var kafkaHost = flag.String("kafkaHost", "localhost:9092", "The Kafka host to consume messages from")
	var logfileName = flag.String("logfile", "", "The file to log messages to")
	var printVersion = flag.Bool("version", false, "Print Holmes version and exit")
	var anonIP = flag.Bool("anonIP", false, "Sets the last octet (IPv4) or the last 80 bits (IPv6) of the client's IP address to 0 in the tracking object before submitting it to Kafka.")
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		return
	}

	initLogging(logfileName)

	var baseURL string

	//only prepend the protocol when the user provided one
	if *protocol != "" {
		baseURL = *protocol + "://"
	}

	//only prepend the host when the user provided one
	if *host != "" {
		baseURL = baseURL + *host
	}

	//we only need to specify the port number when it's different than the
	//standard 80/443.
	if *protocol != "" && !(*protocol == "https" && *proxyPort == "443") && !(*protocol == "http" && *proxyPort == "80") {
		baseURL = baseURL + ":" + *proxyPort
	}
	baseURL = baseURL + *proxyPath

	trackingChannel := make(chan *tracker.TrackingObject, 10)
	trackingParams := tracker.TrackingParams{
		TrackingChannel: trackingChannel,
		AnonymizeIP:     *anonIP,
	}

	http.HandleFunc("/track", provideTrackingParams(trackingParams, tracker.Track))
	http.HandleFunc("/analytics.js", analytics.Analytics(baseURL))
	go publisher.Publish(trackingChannel, kafkaHost, "tracking")

	startServer(*listenHost, *listenPort)
}
