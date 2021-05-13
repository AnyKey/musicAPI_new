package music

import (
	"context"
	"musicAPI/model"
)

type PostgresRepository interface {
	GetTracks(track string, artist string) ([]model.TrackSelect, error)
	SetTracks(newTracks OwnTrack) error
	GetGenreTracks(genre string) ([]model.TrackSelect, error)
	GetArtistTracks(artist string) ([]model.TrackSelect, error)
	GetChart(sortTo string) ([]ChartSelect, error)
}

type RedisRepository interface {
	GetTracksRedis(ctx context.Context, track string, artist string) []model.TrackSelect
	SetTracksRedis(ctx context.Context, track string, artist string, bytes []byte)
	GetGenreRedis(ctx context.Context, genre string) []model.TrackSelect
	SetGenreRedis(ctx context.Context, genre string, bytes []byte)
	GetArtistRedis(ctx context.Context, artist string) []model.TrackSelect
	SetArtistRedis(ctx context.Context, artist string, bytes []byte)
	GetChartRedis(ctx context.Context, sortTo string) []ChartSelect
	SetChartRedis(ctx context.Context, sortTo string, bytes []byte)
	GetAlbumRedis(ctx context.Context, album string, artist string) *Root
	SetAlbumRedis(ctx context.Context, album string, artist string, bytes []byte)
}
type Delivery interface {
	AlbumInfoReq(album string, artist string) (*Root, error)
	TrackSearchReq(track string, artist string) (*OwnTrack, error)
}
type UseCase interface {
	ArtistReq(ctx context.Context, artist string) ([]model.TrackSelect, error)
	GenreReq(ctx context.Context, genre string) ([]model.TrackSelect, error)
	AlbumInfoRes(ctx context.Context, album string, artist string) (*Root, error)
	ChartReq(ctx context.Context, sortTo string) ([]ChartSelect, error)
	TrackReq(ctx context.Context, track string, artist string) ([]model.TrackSelect, bool, error)
}
type ElasticRepository interface {
	ElasticAdd(tracks []model.TrackSelect) error
	ElasticGet(tracks []model.TrackSelect) bool
}

type ChartSelect struct {
	Track     string `json:"track"`
	Artist    string `json:"artist"`
	Album     string `json:"album"`
	Listeners int    `json:"listeners"`
	Playcount int    `json:"playcount"`
	Genre     string `json:"genre"`
}

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
