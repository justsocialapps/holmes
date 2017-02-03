package analytics

import (
	"net/http"
	"strings"

	"github.com/justsocialapps/holmes/assets"
	"github.com/satori/go.uuid"
)

func Analytics(baseUrl string) http.HandlerFunc {
	res := strings.Replace(assets.Analyticsjs, "__HOLMES_BASE_URL__", baseUrl, -1)

	return func(w http.ResponseWriter, r *http.Request) {
		uniqueRes := strings.Replace(res, "__HOLMES_ID__", uuid.NewV4().String(), -1)
		w.Header().Add("Content-Type", "application/javascript")
		w.Write([]byte(uniqueRes))
	}
}
