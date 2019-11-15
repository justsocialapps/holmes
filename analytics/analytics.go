package analytics

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"strings"

	"github.com/justsocialapps/holmes/assets"
)

var analyticsRes []byte
var etag string

// Analytics returns an HTTP handler function that delivers the tracking client
// library.
func Analytics(w http.ResponseWriter, r *http.Request) {
	if analyticsRes == nil || len(analyticsRes) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Add("Content-Type", "application/javascript")
	w.Header().Set("Etag", etag)

	if _, err := w.Write(analyticsRes); err != nil {
		log.Printf("Error sending analytics script: %s\n", err)
	}
}

func PrepareAnalytics(baseURL string) {
	if len(baseURL) != 0 && (analyticsRes == nil || len(analyticsRes) == 0) {
		analyticsRes = []byte(strings.Replace(assets.Analyticsjs, "__HOLMES_BASE_URL__", baseURL, -1))
		hash := sha256.Sum256(analyticsRes)
		etag = hex.EncodeToString(hash[:32])
	}
}
