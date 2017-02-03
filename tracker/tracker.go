package tracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type TrackingObject struct {
	UserAgent string                 `json:"userAgent"`
	Referer   string                 `json:"referer"`
	IPAddress string                 `json:"ipAddress"`
	Time      int64                  `json:"time"`
	Target    map[string]interface{} `json:"target"`
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

	trackingObject := &TrackingObject{
		UserAgent: r.UserAgent(),
		Referer:   r.Referer(),
		IPAddress: strings.Split(r.RemoteAddr, ":")[0],
		Time:      time.Now().Unix(),
		Target:    target,
	}
	log.Println(fmt.Sprintf("tracking object %v", trackingObject))

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
