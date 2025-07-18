package model

type Location struct {
	Username string `json:"username"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64 `json:"timestamp"`
	TimestampFormatted string `json:"timestamp_formatted,omitempty"`
	Triggered bool `json:"triggered"`
	IsRealtime bool `json:"is_realtime"`
}
type LocationResponse struct {
	Locations []Location `json:"locations"`
}
