package main

import (
	"flag"

	"github.com/justsocialapps/holmes/analytics"
)

var (
	protocol     string
	host         string
	proxyPort    string
	proxyPath    string
	listenHost   string
	listenPort   string
	kafkaHost    string
	logfileName  string
	printVersion bool
	anomyzeIP    bool
)

func init() {
	flag.StringVar(&protocol, "protocol", "https", "The protocol used to serve Holmes resources ('http' or 'https')")
	flag.StringVar(&host, "host", "localhost", "The host name used to reach Holmes")
	flag.StringVar(&proxyPort, "proxyPort", "3001", "The TCP port for reaching Holmes if Holmes is operated behind a reverse proxy.")
	flag.StringVar(&proxyPath, "proxyPath", "", "The base path for reaching Holmes if Holmes is operated behind a reverse proxy.")
	flag.StringVar(&listenHost, "listenHost", "", "The host Holmes listens on.")
	flag.StringVar(&listenPort, "listenPort", "3001", "The TCP port that Holmes listens on")
	flag.StringVar(&kafkaHost, "kafkaHost", "localhost:9092", "The Kafka host to consume messages from")
	flag.StringVar(&logfileName, "logfile", "", "The file to log messages to")
	flag.BoolVar(&printVersion, "version", false, "Print Holmes version and exit")
	flag.BoolVar(&anomyzeIP, "anonIP", false, "Sets the last octet (IPv4) or the last 80 bits (IPv6) of the client's IP address to 0 in the tracking object before submitting it to Kafka.")

	flag.Parse()

	var baseURL string

	//only prepend the protocol when the user provided one
	if protocol != "" {
		baseURL = protocol + "://"
	}

	//only prepend the host when the user provided one
	if host != "" {
		baseURL = baseURL + host
	}

	//we only need to specify the port number when it's different than the
	//standard 80/443.
	if protocol != "" && !(protocol == "https" && proxyPort == "443") && !(protocol == "http" && proxyPort == "80") {
		baseURL = baseURL + ":" + proxyPort
	}

	baseURL = baseURL + proxyPath

	analytics.PrepareAnalytics(baseURL)
}
