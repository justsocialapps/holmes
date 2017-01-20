package models

type TrackingObject struct {
	UserAgent string                 `json:"userAgent"`
	Referer   string                 `json:"referer"`
	IPAddress string                 `json:"ipAddress"`
	Time      int64                  `json:"time"`
	Target    map[string]interface{} `json:"target"`
}
