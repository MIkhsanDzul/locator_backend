package model

type Location struct {
	Username string `json:"username"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Triggered bool `json:"triggered"`
	IsRealtime bool `json:"is_realtime"`
}

type LocationResponse struct {
	Locations []Location `json:"locations"`
}

func (l *Location) SetUsernameFromDB(username string) {
	l.Username = username
}
