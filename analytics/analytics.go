package analytics

import (
	"log"
	"net/http"
	"strings"

	"github.com/justsocialapps/holmes/assets"
)

// Analytics returns an HTTP handler function that delivers the tracking client
// library.
func Analytics(baseURL string) http.HandlerFunc {
	res := strings.Replace(assets.Analyticsjs, "__HOLMES_BASE_URL__", baseURL, -1)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/javascript")
		_, err := w.Write([]byte(res))
		if err != nil {
			log.Printf("Error sending analytics script: %s\n", err)
		}
	}
}
