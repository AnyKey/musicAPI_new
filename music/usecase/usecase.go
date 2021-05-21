package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"musicAPI/music"
)

type musicUseCase struct {
	RedisRepo    music.RedisRepository
	PostgresRepo music.PostgresRepository
	ApiRepo      music.ApiDelivery
	EsRepo       music.ElasticDelivery
	GrpcConn     music.GrpcDelivery
}

func New(
	musicRedisRepo music.RedisRepository,
	musicPostgresRepo music.PostgresRepository,
	musicApiRepo music.ApiDelivery,
	musicEsRepo music.ElasticDelivery,
	musicGrpc music.GrpcDelivery,
) music.UseCase {
	return &musicUseCase{
		RedisRepo:    musicRedisRepo,
		PostgresRepo: musicPostgresRepo,
		ApiRepo:      musicApiRepo,
		EsRepo:       musicEsRepo,
		GrpcConn:     musicGrpc,
	}
}

func (muc *musicUseCase) AlbumInfoRes(ctx context.Context, album string, artist string) (*music.Root, error) {
	var err error
	result := muc.RedisRepo.GetAlbumRedis(ctx, album, artist)
	if result != nil {
		return result, nil
		if err != nil {
			return nil, err
		}
	}

	re, err := muc.ApiRepo.AlbumInfoReq(album, artist)
	bytes, err := json.Marshal(re)
	if err == nil {
		muc.RedisRepo.SetAlbumRedis(ctx, album, artist, bytes)
		return re, nil
	}
	return nil, err
}
func (muc *musicUseCase) ChartReq(ctx context.Context, sortTo string) ([]music.ChartSelect, error) {
	chart := muc.RedisRepo.GetChartRedis(ctx, sortTo)
	if chart != nil {
		return chart, nil
	}
	chart, err := muc.PostgresRepo.GetChart(sortTo)
	if err != nil {
		return nil, err
	}
	if chart != nil {
		bytes, err := json.Marshal(chart)
		if err == nil {
			muc.RedisRepo.SetChartRedis(ctx, sortTo, bytes)
		}
		return chart, nil
	}
	return nil, nil
}

func (muc *musicUseCase) ArtistReq(ctx context.Context, artist string) ([]music.TrackSelect, error) {
	tracks := muc.RedisRepo.GetArtistRedis(ctx, artist)
	if tracks != nil {
		return tracks, nil
	}
	tracks, err := muc.PostgresRepo.GetArtistTracks(artist)
	if err != nil {
		return nil, err
	}
	if tracks != nil {
		bytes, err := json.Marshal(tracks)
		if err != nil {
			log.Println(err)
		}
		if err == nil {
			muc.RedisRepo.SetArtistRedis(ctx, artist, bytes)
		}
		return tracks, nil
	}
	return nil, nil
}

func (muc *musicUseCase) GenreReq(ctx context.Context, genre string) ([]music.TrackSelect, error) {
	tracks := muc.RedisRepo.GetGenreRedis(ctx, genre)
	if tracks != nil {
		return tracks, nil
	}
	tracks, err := muc.PostgresRepo.GetGenreTracks(genre)
	if err != nil {
		return nil, err
	}
	if tracks != nil {
		bytes, err := json.Marshal(tracks)
		if err != nil {
			log.Println(err)
		}
		if err == nil {
			muc.RedisRepo.SetGenreRedis(ctx, genre, bytes)
		}
		return tracks, nil
	}
	return nil, nil
}

func (muc *musicUseCase) TrackReq(ctx context.Context, track string, artist string) ([]music.TrackSelect, bool, error) {
	var value bool
	var result []music.TrackSelect
	tracks := muc.RedisRepo.GetTracksRedis(ctx, track, artist)
	if tracks != nil {
		if muc.EsRepo.ElasticGet(tracks) != true {
			err := muc.EsRepo.ElasticAdd(tracks)
			if err != nil {
				log.Println("error elastic add")
			}
		}
	}
	if tracks != nil {
		value = true
		return tracks, value, nil
	}

	tracks, err := muc.PostgresRepo.GetTracks(track, artist)
	bytes, err := json.Marshal(tracks)
	if err == nil && tracks != nil {
		muc.RedisRepo.SetTracksRedis(ctx, track, artist, bytes)
		value = true
	}
	if tracks != nil {
		if tracks[0].Album == "" || tracks[0].Name == "" || tracks[0].Artist == "" {
			value = false
		}
		if tracks[0].Album != "" && tracks[0].Name != "" && tracks[0].Artist != "" {
			return tracks, value, nil
		}
	}
	if tracks == nil {
		re, err := muc.ApiRepo.TrackSearchReq(track, artist)
		if err != nil {
			fmt.Println(err.Error())
		}
		if re == nil {
			value = false
			return nil, value, err
		}
		if re != nil {
			go func() {
				err = muc.PostgresRepo.SetTracks(*re)
				if err != nil {
					fmt.Println(err.Error())
				}
			}()
			value = true
		}
		result = structConv(re)

	}
	if result == nil {
		if tracks[0].Album == "" || tracks[0].Name == "" || tracks[0].Artist == "" {
			value = false
			return nil, value, err
		}
	}
	return result, value, nil
}

func (muc *musicUseCase) SetLike(name string, artist string, token string) (*string, error) {
	check, err := muc.PostgresRepo.CheckTrack(name, artist)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, errors.New("Track doesnt exist")
	}
	message, err := muc.GrpcConn.SetLike(name, artist, token)
	if err != nil {
		return nil, err
	}
	return message, nil
}
func (muc *musicUseCase) GetLike(name string, artist string, token string) (*music.LikeList, error) {
	check, err := muc.PostgresRepo.CheckTrack(name, artist)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, errors.New("Track doesnt exist")
	}
	message, err := muc.GrpcConn.GetLike(name, artist, token)
	if err != nil {
		return nil, err
	}
	return message, nil
}
func structConv(trackList *music.OwnTrack) []music.TrackSelect {
	return []music.TrackSelect{{
		trackList.Name,
		trackList.Album.Artist,
		trackList.Album.Album,
	},
	}
}
