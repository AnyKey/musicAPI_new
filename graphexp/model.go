package graphexp

import "musicAPI/graph/model"

type Repository interface {
	GetTracks() ([]*model.Tracks, error)
}
type UseCase interface {
	GetTracks() ([]*model.Tracks, error)
}
type Tracks struct {
	Name   string `db:"track" json:"name"`
	Artist string `db:"artist" json:"artist"`
	Album  string `db:"album" json:"album"`
}
