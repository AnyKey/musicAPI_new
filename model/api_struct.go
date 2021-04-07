package model

//Album
type Track struct {
	Name     string `json:"name"`
	Duration string `json:"duration"`
}

type Album struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Tracks struct {
		Tracks []Track `json:"track"`
	} `json:"tracks"`
}

type Root struct {
	Album Album `json:"album"`
}

//Tracks
type TrackSearch struct {
	Name      string `json:"name"`
	Artist    string `json:"artist"`
	Listeners string `json:"listeners"`
}
type ResTrackSearch struct {
	Result struct {
		Trackmatches struct {
			Track []TrackSearch `json:"track"`
		} `json:"trackmatches"`
	} `json:"results"`
}
type OwnTrack struct {
	Name      string     `json:"name"`
	Album     TrackAlbum `json:"album"`
	Listeners string     `json:"listeners"`
	Playcount string     `json:"playcount"`
	TopTags   struct {
		Genre []Tags `json:"tag"`
	} `json:"toptags"`
}
type Tags struct {
	Tag string `json:"name"`
}
type TrackAlbum struct {
	Artist string `json:"artist"`
	Album  string `json:"title"`
}
type TrackRoot struct {
	Track OwnTrack `json:"track"`
}
