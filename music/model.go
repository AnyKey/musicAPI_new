package music

import (
	"context"
)

type PostgresRepository interface {
	GetTracks(string, string) ([]TrackSelect, error)
	SetTracks(OwnTrack) error
	GetGenreTracks(string) ([]TrackSelect, error)
	GetArtistTracks(string) ([]TrackSelect, error)
	GetChart(string) ([]ChartSelect, error)
}

type RedisRepository interface {
	GetTracksRedis(context.Context, string, string) []TrackSelect
	SetTracksRedis(context.Context, string, string, []byte)
	GetGenreRedis(context.Context, string) []TrackSelect
	SetGenreRedis(context.Context, string, []byte)
	GetArtistRedis(context.Context, string) []TrackSelect
	SetArtistRedis(context.Context, string, []byte)
	GetChartRedis(context.Context, string) []ChartSelect
	SetChartRedis(context.Context, string, []byte)
	GetAlbumRedis(context.Context, string, string) *Root
	SetAlbumRedis(context.Context, string, string, []byte)
}
type ApiDelivery interface {
	AlbumInfoReq(string, string) (*Root, error)
	TrackSearchReq(string, string) (*OwnTrack, error)
}
type UseCase interface {
	ArtistReq(context.Context, string) ([]TrackSelect, error)
	GenreReq(context.Context, string) ([]TrackSelect, error)
	AlbumInfoRes(context.Context, string, string) (*Root, error)
	ChartReq(context.Context, string) ([]ChartSelect, error)
	TrackReq(context.Context, string, string) ([]TrackSelect, bool, error)
}
type ElasticDelivery interface {
	ElasticAdd([]TrackSelect) error
	ElasticGet([]TrackSelect) bool
}
type TrackSelect struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
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
