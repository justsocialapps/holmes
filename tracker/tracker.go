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

type UserAgent struct {
	Bot            bool   `json:"bot"`
	Mobile         bool   `json:"mobile"`
	BrowserName    string `json:"browserName"`
	BrowserVersion string `json:"browserVersion"`
	Locale         string `json:"locale"`
	OS             string `json:"os"`
	Platform       string `json:"platform"`
}

type TrackingObject struct {
	UA        UserAgent              `json:"ua"`
	UserAgent string                 `json:"userAgent"`
	Referer   string                 `json:"referer"`
	IPAddress string                 `json:"ipAddress"`
	Time      int64                  `json:"time"`
	Target    map[string]interface{} `json:"target"`
}

func anonymizeIP(ip string) (string) {
	ipAddress := net.ParseIP(ip)
        if ipAddress.To4() == nil { 
	       // mask IPv6 address
               return ipAddress.Mask(net.CIDRMask(48, 128)).String()
        } else {
               return ipAddress.Mask(net.CIDRMask(24, 32)).String()
        }
}


func composeTrackingObject(r *http.Request) (*TrackingObject, error) {
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
	anonymizedIpAddress := anonymizeIP(originIPAddress)

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
		IPAddress: anonymizedIpAddress,
		Time:      time.Now().Unix(),
		Target:    target,
	}

	return trackingObject, nil
}

func Track(out chan<- *TrackingObject, w http.ResponseWriter, r *http.Request) {
	trackingObject, err := composeTrackingObject(r)
	if err != nil {
		log.Printf("Error processing tracking request: %s\n", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	out <- trackingObject
	w.WriteHeader(http.StatusNoContent)
}
