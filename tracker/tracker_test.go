package tracker

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justsocialapps/assert"
)

func TestTrackWithNoRequestArgumentDoesNotPublishTrackingEvent(t *testing.T) {
	assert := assert.NewAssert(t)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "http://example.com", nil)
	trackingChannel := make(chan *TrackingObject)
	done := make(chan struct{})
	go func() {
		Track(TrackingParams{trackingChannel,false,}, recorder, request)
		close(done)
	}()
	select {
	case trackingObject := <-trackingChannel:
		t.Errorf("Expected no tracking object but got one: %v", trackingObject)
	case <-done:
		assert.Equal(recorder.Code, http.StatusNoContent, "Unexpected status code")
	}
}
func TestTrackWithWrongRequestDoesNotPublishTrackingEvent(t *testing.T) {
	assert := assert.NewAssert(t)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "http://example.com/?t=123", nil)
	trackingChannel := make(chan *TrackingObject)
	done := make(chan struct{})
	go func() {
		Track(TrackingParams{trackingChannel, false,}, recorder, request)
		close(done)
	}()
	select {
	case trackingObject := <-trackingChannel:
		t.Errorf("Expected no tracking object but got one: %v", trackingObject)
	case <-done:
		assert.Equal(recorder.Code, http.StatusNoContent, "Unexpected status code")
	}
}

func TestTrackWithCorrectRequestPublishesTrackingEvent(t *testing.T) {
	assert := assert.NewAssert(t)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "http://example.org?t=%7B%22hash%22%3A%22entity.865%22%2C%22entity%22%3A%22ENTITY%2C865%22%2C%22pageTitle%22%3A%22Just%20Software%20(Organisation)%22%2C%22type%22%3A%22PAGE_VIEW%22%2C%22holmesId%22%3A%227be4c968-aea0-4d76-a534-49bd0bb0222b%22%7D", nil)
	request.Header.Set("User-Agent", "go test")
	request.Header.Set("Referer", "referer")
	request.Header.Set("X-Forwarded-For", "1.2.3.4")

	trackingChannel := make(chan *TrackingObject)
	done := make(chan struct{})
	go func() {
		Track(TrackingParams{trackingChannel, false,}, recorder, request)
		close(done)
	}()
	select {
	case trackingObject := <-trackingChannel:
		if trackingObject.Time == 0 {
			t.Errorf("tracking object's time is 0")
		}
		assert.Equal(trackingObject.UserAgent, "go test", "wrong user agent")
		assert.Equal(trackingObject.Referer, "referer", "wrong referer")
		assert.Equal(trackingObject.Target["hash"], "entity.865", "wrong hash")
		assert.Equal(trackingObject.Target["entity"], "ENTITY,865", "wrong entity")
		assert.Equal(trackingObject.Target["pageTitle"], "Just Software (Organisation)", "wrong page title")
		assert.Equal(trackingObject.Target["type"], "PAGE_VIEW", "wrong type")
		assert.Equal(trackingObject.Target["holmesId"], "7be4c968-aea0-4d76-a534-49bd0bb0222b", "wrong holmes id")
		assert.Equal(trackingObject.IPAddress, "1.2.3.4", "wrong IP address")
	case <-done:
		assert.Equal(recorder.Code, http.StatusNoContent, "Unexpected status code")
	}
}

func TestIPv4IsAnonymizedInTrackingObject(t *testing.T) {
	assert := assert.NewAssert(t)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "http://example.org?t=%7B%22hash%22%3A%22entity.865%22%2C%22entity%22%3A%22ENTITY%2C865%22%2C%22pageTitle%22%3A%22Just%20Software%20(Organisation)%22%2C%22type%22%3A%22PAGE_VIEW%22%2C%22holmesId%22%3A%227be4c968-aea0-4d76-a534-49bd0bb0222b%22%7D", nil)
	request.Header.Set("User-Agent", "go test")
	request.Header.Set("Referer", "referer")
	request.Header.Set("X-Forwarded-For", "1.2.3.4")

	trackingChannel := make(chan *TrackingObject)
	done := make(chan struct{})
	go func() {
		Track(TrackingParams{trackingChannel, true,}, recorder, request)
		close(done)
	}()
	select {
	case trackingObject := <-trackingChannel:
		assert.Equal(trackingObject.IPAddress, "1.2.3.0", "wrong IP address")
	case <-done:
		assert.Equal(recorder.Code, http.StatusNoContent, "Unexpected status code")
	}
}

func TestIPv6IsAnonymizedInTrackingObject(t *testing.T) {
	assert := assert.NewAssert(t)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "http://example.org?t=%7B%22hash%22%3A%22entity.865%22%2C%22entity%22%3A%22ENTITY%2C865%22%2C%22pageTitle%22%3A%22Just%20Software%20(Organisation)%22%2C%22type%22%3A%22PAGE_VIEW%22%2C%22holmesId%22%3A%227be4c968-aea0-4d76-a534-49bd0bb0222b%22%7D", nil)
	request.Header.Set("User-Agent", "go test")
	request.Header.Set("Referer", "referer")
	request.Header.Set("X-Forwarded-For", "2001:0db8:85a3:0000:0000:8a2e:0370:7334")

	trackingChannel := make(chan *TrackingObject)
	done := make(chan struct{})
	go func() {
		Track(TrackingParams{trackingChannel, true,}, recorder, request)
		close(done)
	}()
	select {
	case trackingObject := <-trackingChannel:
		assert.Equal(trackingObject.IPAddress, "2001:db8:85a3::", "wrong IP address")
	case <-done:
		assert.Equal(recorder.Code, http.StatusNoContent, "Unexpected status code")
	}
}


func TestIPv6LocalhostIsAnonymizedInTrackingObject(t *testing.T) {
	assert := assert.NewAssert(t)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "http://example.org?t=%7B%22hash%22%3A%22entity.865%22%2C%22entity%22%3A%22ENTITY%2C865%22%2C%22pageTitle%22%3A%22Just%20Software%20(Organisation)%22%2C%22type%22%3A%22PAGE_VIEW%22%2C%22holmesId%22%3A%227be4c968-aea0-4d76-a534-49bd0bb0222b%22%7D", nil)
	request.Header.Set("User-Agent", "go test")
	request.Header.Set("Referer", "referer")
	request.Header.Set("X-Forwarded-For", "::1")

	trackingChannel := make(chan *TrackingObject)
	done := make(chan struct{})
	go func() {
		Track(TrackingParams{trackingChannel, true,}, recorder, request)
		close(done)
	}()
	select {
	case trackingObject := <-trackingChannel:
		assert.Equal(trackingObject.IPAddress, "::", "wrong IP address")
	case <-done:
		assert.Equal(recorder.Code, http.StatusNoContent, "Unexpected status code")
	}
}
