package models

type TrackingTarget struct {
	Hash      string `json:"hash"`
	Entity    string `json:"entity"`
	PageTitle string `json:"pageTitle"`
	Type      string `json:"type"`
	HolmesId  string `json:"holmesId"`
}

type TrackingObject struct {
	UserAgent string         `json:"userAgent"`
	Referer   string         `json:"referer"`
	IPAddress string         `json:"ipAddress"`
	Time      int64          `json:"time"`
	Target    TrackingTarget `json:"target"`
}
