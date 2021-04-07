package model

type TrackSelect struct {
	Name   string `db:"track" json:"name"`
	Artist string `db:"artist" json:"artist"`
	Album  string `db:"album" json:"album"`
}
type ChartSelect struct {
	Track     string
	Artist    string
	Album     string
	Listeners int
	Playcount int
	Genre     string
}
