package model

type Location struct {
	Username string `json:"username"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64 `json:"timestamp"`
	TimestampFormatted string `json:"timestamp_formatted,omitempty"`
	Triggered bool `json:"triggered"`
}
type LocationResponse struct {
	Locations []Location `json:"locations"`
}
