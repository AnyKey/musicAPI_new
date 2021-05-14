package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"musicAPI/music"
	"time"
)

type Repository struct {
	Redis *redis.Client
}

func New(redis *redis.Client) *Repository {
	return &Repository{
		Redis: redis,
	}
}

func (repo *Repository) GetTracksRedis(ctx context.Context, track string, artist string) []music.TrackSelect {
	var trackList []music.TrackSelect
	res := repo.Redis.Get(ctx, "Track:"+track+"_Artist:"+artist)
	if res.Err() != nil {
		return nil
	}
	bytes, err := res.Bytes()
	if err != nil {
		return nil
	}
	err = json.Unmarshal(bytes, &trackList)
	if err != nil {
		return nil
	}
	return trackList
}

func (repo *Repository) GetGenreRedis(ctx context.Context, genre string) []music.TrackSelect {
	var trackList []music.TrackSelect
	res := repo.Redis.Get(ctx, "Genre:"+genre)
	if res.Err() != nil {
		return nil
	}
	bytes, err := res.Bytes()
	if err != nil {
		return nil
	}
	err = json.Unmarshal(bytes, &trackList)
	if err != nil {
		return nil
	}
	return trackList
}

func (repo *Repository) GetArtistRedis(ctx context.Context, artist string) []music.TrackSelect {
	var trackList []music.TrackSelect
	res := repo.Redis.Get(ctx, "Artist:"+artist)
	if res.Err() != nil {
		return nil
	}
	bytes, err := res.Bytes()
	if err != nil {
		return nil
	}
	err = json.Unmarshal(bytes, &trackList)
	if err != nil {
		return nil
	}
	return trackList
}

func (repo *Repository) GetChartRedis(ctx context.Context, sortTo string) []music.ChartSelect {
	var trackList []music.ChartSelect
	res := repo.Redis.Get(ctx, "SortTo:"+sortTo)
	if res.Err() != nil {
		return nil
	}
	bytes, err := res.Bytes()
	if err != nil {
		return nil
	}
	err = json.Unmarshal(bytes, &trackList)
	if err != nil {
		return nil
	}
	return trackList
}

func (repo *Repository) GetAlbumRedis(ctx context.Context, album string, artist string) *music.Root {
	var trackList *music.Root
	res := repo.Redis.Get(ctx, "Album:"+album+"_Artist:"+artist)
	if res.Err() != nil {
		return nil
	}
	bytes, err := res.Bytes()
	if err != nil {
		return nil
	}
	err = json.Unmarshal(bytes, &trackList)
	if err != nil {
		return nil
	}
	return trackList
}

func (repo *Repository) SetAlbumRedis(ctx context.Context, album string, artist string, bytes []byte) {
	repo.Redis.Set(ctx, "Album:"+album+"_Artist:"+artist, bytes, 5*time.Minute)
	return
}
func (repo *Repository) SetChartRedis(ctx context.Context, sortTo string, bytes []byte) {
	repo.Redis.Set(ctx, "SortTo:"+sortTo, bytes, 5*time.Minute)
	return
}
func (repo *Repository) SetArtistRedis(ctx context.Context, artist string, bytes []byte) {
	repo.Redis.Set(ctx, "Artist:"+artist, bytes, 5*time.Minute)
	return
}
func (repo *Repository) SetGenreRedis(ctx context.Context, genre string, bytes []byte) {
	repo.Redis.Set(ctx, "Genre:"+genre, bytes, 5*time.Minute)
	return
}
func (repo *Repository) SetTracksRedis(ctx context.Context, track string, artist string, bytes []byte) {
	repo.Redis.Set(ctx, "Track:"+track+"_Artist:"+artist, bytes, 20*time.Minute)
	return
}
