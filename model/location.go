package model

type Location struct {
	Name string `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64 `json:"timestamp"`
}