package analytics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justsocialapps/assert"
	"github.com/justsocialapps/holmes/assets"
)

func TestAnalytics(t *testing.T) {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "", nil)
	assets.Analyticsjs = "base url: __HOLMES_BASE_URL__"
	assert := assert.NewAssert(t)

	Analytics("https://example.org/baseurl")(recorder, request)

	assert.Equal(recorder.Header().Get("Content-Type"), "application/javascript", "wrong content type")
	assert.Equal(recorder.Code, http.StatusOK, "wrong status code")
	assert.Match("base url: https://example.org/baseurl", recorder.Body.String(), "wrong body")
}
