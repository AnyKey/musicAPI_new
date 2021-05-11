package music

import (
	"context"
	"musicAPI/model"
)

const (
	BaseURL = "http://ws.audioscrobbler.com/2.0/"
	ApiKey  = "d84296d9388306355db600e324a85b9b"
)

type PostgresRepository interface {
	GetTracks(track string, artist string) ([]model.TrackSelect, error)
	SetTracks(newTracks model.OwnTrack) error
	GetGenreTracks(genre string) ([]model.TrackSelect, error)
	GetArtistTracks(artist string) ([]model.TrackSelect, error)
	GetChart(sortTo string) ([]model.ChartSelect, error)
}

type RedisRepository interface {
	GetTracksRedis(ctx context.Context, track string, artist string) []model.TrackSelect
	SetTracksRedis(ctx context.Context, track string, artist string, bytes []byte)
	GetGenreRedis(ctx context.Context, genre string) []model.TrackSelect
	SetGenreRedis(ctx context.Context, genre string, bytes []byte)
	GetArtistRedis(ctx context.Context, artist string) []model.TrackSelect
	SetArtistRedis(ctx context.Context, artist string, bytes []byte)
	GetChartRedis(ctx context.Context, sortTo string) []model.ChartSelect
	SetChartRedis(ctx context.Context, sortTo string, bytes []byte)
	GetAlbumRedis(ctx context.Context, album string, artist string) *model.Root
	SetAlbumRedis(ctx context.Context, album string, artist string, bytes []byte)
}
type ApiRepository interface {
	AlbumInfoReq(album string, artist string) (*model.Root, error)
	TrackSearchReq(track string, artist string) (*model.OwnTrack, error)
}
type UseCase interface {
	ArtistReq(ctx context.Context, artist string) ([]model.TrackSelect, error)
	GenreReq(ctx context.Context, genre string) ([]model.TrackSelect, error)
	AlbumInfoRes(ctx context.Context, album string, artist string) (*model.Root, error)
	ChartReq(ctx context.Context, sortTo string) ([]model.ChartSelect, error)
	TrackReq(ctx context.Context, track string, artist string) ([]model.TrackSelect, bool, error) //very hard!!!
}
type ElasticRepository interface {
	ElasticAdd(tracks []model.TrackSelect) error
	ElasticGet(tracks []model.TrackSelect) bool
}
