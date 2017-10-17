package tracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	ua "github.com/mssola/user_agent"
)

// UserAgent holds values that are provided by the browser's "User-Agent" HTTP
// header.
type UserAgent struct {
	Bot            bool   `json:"bot"`
	Mobile         bool   `json:"mobile"`
	BrowserName    string `json:"browserName"`
	BrowserVersion string `json:"browserVersion"`
	Locale         string `json:"locale"`
	OS             string `json:"os"`
	Platform       string `json:"platform"`
}

// TrackingObject holds all data that is eventually sent to Kafka.
type TrackingObject struct {
	UA        UserAgent              `json:"ua"`
	UserAgent string                 `json:"userAgent"`
	Referer   string                 `json:"referer"`
	IPAddress string                 `json:"ipAddress"`
	Time      int64                  `json:"time"`
	Target    map[string]interface{} `json:"target"`
}

// TrackingParams encapsulates the parameters passed to Track.
type TrackingParams struct {
	TrackingChannel chan<- *TrackingObject
	AnonymizeIP     bool
}

func anonymizeIP(ip string) string {
	ipAddress := net.ParseIP(ip)
	if ipAddress.To4() == nil {
		// mask IPv6 address
		return ipAddress.Mask(net.CIDRMask(48, 128)).String()
	}
	return ipAddress.Mask(net.CIDRMask(24, 32)).String()
}

func composeTrackingObject(anonIP bool, r *http.Request) (*TrackingObject, error) {
	query := r.URL.Query()
	rawTarget := query["t"]
	if len(rawTarget) == 0 {
		return nil, errors.New("No object to track")
	}

	var target map[string]interface{}
	err := json.Unmarshal([]byte(rawTarget[0]), &target)
	if err != nil {
		return nil, fmt.Errorf("Error parsing tracking target '%s'", rawTarget)
	}

	// We prefer the IP address stated in the X-Forwarded-For HTTP header.
	// Only if this header is empty we use http.Request.RemoteAddr.
	var originIPAddress string
	forwardedFor := strings.TrimSpace(strings.SplitN(r.Header.Get("x-forwarded-for"), ",", 2)[0])
	if len(forwardedFor) > 0 {
		originIPAddress = forwardedFor
	} else {
		originIPAddress = strings.Split(r.RemoteAddr, ":")[0]
	}
	var trackingIPAddress string
	if anonIP {
		trackingIPAddress = anonymizeIP(originIPAddress)
	} else {
		trackingIPAddress = originIPAddress
	}

	userAgent := ua.New(r.UserAgent())
	browserName, browserVersion := userAgent.Browser()

	trackingObject := &TrackingObject{
		UA: UserAgent{
			Bot:            userAgent.Bot(),
			Mobile:         userAgent.Mobile(),
			BrowserName:    browserName,
			BrowserVersion: browserVersion,
			Locale:         userAgent.Localization(),
			OS:             userAgent.OS(),
			Platform:       userAgent.Platform(),
		},
		UserAgent: userAgent.UA(),
		Referer:   r.Referer(),
		IPAddress: trackingIPAddress,
		Time:      time.Now().Unix(),
		Target:    target,
	}

	return trackingObject, nil
}

// Track is an HTTP handler function that sends tracking objects provided by
// browsers to a channel.
func Track(params TrackingParams, w http.ResponseWriter, r *http.Request) {
	trackingObject, err := composeTrackingObject(params.AnonymizeIP, r)
	if err != nil {
		log.Printf("Error processing tracking request: %s\n", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	params.TrackingChannel <- trackingObject
	w.WriteHeader(http.StatusNoContent)
}
