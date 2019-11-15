package analytics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justsocialapps/assert"
	"github.com/justsocialapps/holmes/assets"
)

func TestAnalyticsWithoutSetValue(t *testing.T) {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "", nil)
	assets.Analyticsjs = "base url: __HOLMES_BASE_URL__"
	assertion := assert.NewAssert(t)

	Analytics(recorder, request)

	assertion.Equal(len(recorder.Header()), 0, "wrong content type")
	assertion.Equal(recorder.Code, http.StatusNotFound, "wrong status code")
	assertion.Match("", recorder.Body.String(), "wrong body")
}

func TestAnalyticsWithSetValue(t *testing.T) {
	etag := "11b9fdd2cfc348e9f1fbf2b774ca1625ea8bf5ed63a6e17fae6c2b0271895192"
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "", nil)
	assets.Analyticsjs = "base url: __HOLMES_BASE_URL__"
	PrepareAnalytics("https://example.org/baseurl")
	assertion := assert.NewAssert(t)

	Analytics(recorder, request)

	assertion.Equal(recorder.Header().Get("Content-Type"), "application/javascript", "wrong content type")
	assertion.Equal(recorder.Header().Get("Etag"), etag, "wrong etag code")
	assertion.Equal(recorder.Code, http.StatusOK, "wrong status code")
	assertion.Match("base url: https://example.org/baseurl", recorder.Body.String(), "wrong body")
}
