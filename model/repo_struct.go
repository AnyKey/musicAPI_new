package model

type TrackSelect struct {
	Name   string `db:"track" json:"name"`
	Artist string `db:"artist" json:"artist"`
	Album  string `db:"album" json:"album"`
}
type ChartSelect struct {
	Track     string `json:"track"`
	Artist    string `json:"artist"`
	Album     string `json:"album"`
	Listeners int    `json:"listeners"`
	Playcount int    `json:"playcount"`
	Genre     string `json:"genre"`
}
