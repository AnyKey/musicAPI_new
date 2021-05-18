package graphexp

import "musicAPI/graph/model"

type Repository interface {
	GetTracks() ([]*model.Tracks, error)
	GetArtist(*string) ([]*model.Tracks, error)
	SetTrack(OwnTrack) error
}
type UseCase interface {
	GetTracks() ([]*model.Tracks, error)
	GetArtist(*string) ([]*model.Tracks, error)
	SetTrack(model.NewTrack) error
}
type Tracks struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}
type OwnTrack struct {
	Name      string `json:"name"`
	Album     string `json:"album"`
	Artist    string `json:"artist"`
	Listeners string `json:"listeners"`
	Playcount string `json:"playcount"`
	Tag       string `json:"tag"`
}
