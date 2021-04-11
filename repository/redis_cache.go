package repository

import (
	"context"
	"encoding/json"
	"musicAPI/model"
	"time"
)

func (repo Repository) GetTracksRedis(track string, artist string) []model.TrackSelect {
	var trackList []model.TrackSelect
	var ctx = context.Background()
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

func (repo Repository) GetGenreRedis(genre string) []model.TrackSelect {
	var trackList []model.TrackSelect
	var ctx = context.Background()
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

func (repo Repository) GetArtistRedis(artist string) []model.TrackSelect {
	var trackList []model.TrackSelect
	var ctx = context.Background()
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

func (repo Repository) GetChartRedis(sortTo string) []model.ChartSelect {
	var trackList []model.ChartSelect
	var ctx = context.Background()
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

func (repo Repository) GetAlbumRedis(album string, artist string) *model.Root {
	var trackList *model.Root
	var ctx = context.Background()
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

func (repo Repository) SetTokens(ctx context.Context, user string, tokens model.Tokens) error {
	bytes, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	res := repo.Redis.Set(ctx, "JWT:"+user, bytes, 1*time.Hour)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (repo Repository) GetTokens(ctx context.Context, user string) *model.Tokens {
	var tokens model.Tokens
	res := repo.Redis.Get(ctx, "JWT:"+user)
	if res.Err() != nil {
		return nil
	}
	bytes, err := res.Bytes()
	if err != nil {
		return nil
	}
	err = json.Unmarshal(bytes, &tokens)
	if err != nil {
		return nil
	}
	return &tokens

}
