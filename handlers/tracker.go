package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/justsocialapps/holmes/models"
)

func composeTrackingObject(r *http.Request) (*models.TrackingObject, error) {
	query := r.URL.Query()
	rawTarget := query["t"]
	if len(rawTarget) == 0 {
		return nil, errors.New("No object to track")
	}
	var target models.TrackingTarget
	err := json.Unmarshal([]byte(rawTarget[0]), &target)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing tracking target '%s'", rawTarget))
	}

	trackingObject := &models.TrackingObject{
		UserAgent: r.UserAgent(),
		Referer:   r.Referer(),
		IPAddress: strings.Split(r.RemoteAddr, ":")[0],
		Time:      time.Now().Unix(),
		Target:    target,
	}
	log.Println(fmt.Sprintf("tracking object %s", trackingObject))

	return trackingObject, nil
}

func Track(out chan<- *models.TrackingObject, w http.ResponseWriter, r *http.Request) {
	trackingObject, err := composeTrackingObject(r)
	if err != nil {
		log.Printf("Error processing tracking request: %s\n", err)
	}

	out <- trackingObject
	w.WriteHeader(http.StatusNoContent)
}
