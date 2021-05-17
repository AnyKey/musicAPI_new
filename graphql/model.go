package graphql

type UseCase interface {
	GetTracks() ([]TrackSelect, error)
}
type Repository interface {
	GetTracks() ([]TrackSelect, error)
}
type Delivery interface {
	GetTracks([]TrackSelect) error
}
type TrackSelect struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}
