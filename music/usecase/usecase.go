package usecase

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"log"
	"musicAPI/model"
	"musicAPI/music"
)

type musicUseCase struct {
	MusicRedisRepo    music.RedisRepository
	MusicPostgresRepo music.PostgresRepository
	MusicApiRepo      music.ApiRepository
	MusicEsRepo       music.ElasticRepository
}

func New(musicRedisRepo music.RedisRepository, musicPostgresRepo music.PostgresRepository, musicApiRepo music.ApiRepository, musicEsRepo music.ElasticRepository) music.UseCase {
	return &musicUseCase{
		MusicRedisRepo:    musicRedisRepo,
		MusicPostgresRepo: musicPostgresRepo,
		MusicApiRepo:      musicApiRepo,
		MusicEsRepo:       musicEsRepo,
	}
}

func (muc *musicUseCase) AlbumInfoRes(ctx context.Context, album string, artist string) (*model.Root, error) {
	var err error
	result := muc.MusicRedisRepo.GetAlbumRedis(ctx, album, artist)
	if result != nil {
		return result, nil
		if err != nil {
			return nil, err
		}
	}

	re, err := muc.MusicApiRepo.AlbumInfoReq(album, artist)
	bytes, err := json.Marshal(re)
	if err == nil {
		muc.MusicRedisRepo.SetAlbumRedis(ctx, album, artist, bytes)
		return re, nil
	}
	return nil, err
}
func (muc *musicUseCase) ChartReq(ctx context.Context, sortTo string) ([]model.ChartSelect, error) {
	chart := muc.MusicRedisRepo.GetChartRedis(ctx, sortTo)
	if chart != nil {
		return chart, nil
	}
	chart, err := muc.MusicPostgresRepo.GetChart(sortTo)
	if err != nil {
		return nil, err
	}
	if chart != nil {
		bytes, err := json.Marshal(chart)
		if err == nil {
			muc.MusicRedisRepo.SetChartRedis(ctx, sortTo, bytes)
		}
		return chart, nil
	}
	return nil, nil
}

func (muc *musicUseCase) ArtistReq(ctx context.Context, artist string) ([]model.TrackSelect, error) {
	tracks := muc.MusicRedisRepo.GetArtistRedis(ctx, artist)
	if tracks != nil {
		return tracks, nil
	}
	tracks, err := muc.MusicPostgresRepo.GetArtistTracks(artist)
	if err != nil {
		return nil, err
	}
	if tracks != nil {
		bytes, err := json.Marshal(tracks)
		if err != nil {
			log.Println(err)
		}
		if err == nil {
			muc.MusicRedisRepo.SetArtistRedis(ctx, artist, bytes)
		}
		return tracks, nil
	}
	return nil, nil
}

func (muc *musicUseCase) GenreReq(ctx context.Context, genre string) ([]model.TrackSelect, error) {
	tracks := muc.MusicRedisRepo.GetGenreRedis(ctx, genre)
	if tracks != nil {
		return tracks, nil
	}
	tracks, err := muc.MusicPostgresRepo.GetGenreTracks(genre)
	if err != nil {
		return nil, err
	}
	if tracks != nil {
		bytes, err := json.Marshal(tracks)
		if err != nil {
			log.Println(err)
		}
		if err == nil {
			muc.MusicRedisRepo.SetGenreRedis(ctx, genre, bytes)
		}
		return tracks, nil
	}
	return nil, nil
}

/*func structConv(trackList *model.OwnTrack) model.TrackSelect {
	return model.TrackSelect{trackList.Name, trackList.Album.Artist, trackList.Album.Album}
}*/
