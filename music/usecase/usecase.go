package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"musicAPI/music"
)

type musicUseCase struct {
	MusicRedisRepo    music.RedisRepository
	MusicPostgresRepo music.PostgresRepository
	MusicApiRepo      music.ApiDelivery
	MusicEsRepo       music.ElasticDelivery
}

func New(musicRedisRepo music.RedisRepository,
	musicPostgresRepo music.PostgresRepository,
	musicApiRepo music.ApiDelivery,
	musicEsRepo music.ElasticDelivery) music.UseCase {
	return &musicUseCase{
		MusicRedisRepo:    musicRedisRepo,
		MusicPostgresRepo: musicPostgresRepo,
		MusicApiRepo:      musicApiRepo,
		MusicEsRepo:       musicEsRepo,
	}
}

func (muc *musicUseCase) AlbumInfoRes(ctx context.Context, album string, artist string) (*music.Root, error) {
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
func (muc *musicUseCase) ChartReq(ctx context.Context, sortTo string) ([]music.ChartSelect, error) {
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

func (muc *musicUseCase) ArtistReq(ctx context.Context, artist string) ([]music.TrackSelect, error) {
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

func (muc *musicUseCase) GenreReq(ctx context.Context, genre string) ([]music.TrackSelect, error) {
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

func (muc *musicUseCase) TrackReq(ctx context.Context, track string, artist string) ([]music.TrackSelect, bool, error) {
	var value bool
	var result []music.TrackSelect
	tracks := muc.MusicRedisRepo.GetTracksRedis(ctx, track, artist)
	if tracks != nil {
		if muc.MusicEsRepo.ElasticGet(tracks) != true {
			err := muc.MusicEsRepo.ElasticAdd(tracks)
			if err != nil {
				log.Println("error elastic add")
			}
		}
	}
	if tracks != nil {
		value = true
		return tracks, value, nil
	}

	tracks, err := muc.MusicPostgresRepo.GetTracks(track, artist)
	bytes, err := json.Marshal(tracks)
	if err == nil && tracks != nil {
		muc.MusicRedisRepo.SetTracksRedis(ctx, track, artist, bytes)
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
		re, err := muc.MusicApiRepo.TrackSearchReq(track, artist)
		if err != nil {
			fmt.Println(err.Error())
		}
		if re == nil {
			value = false
			return nil, value, err
		}
		if re != nil {
			go func() {
				err = muc.MusicPostgresRepo.SetTracks(*re)
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
func structConv(trackList *music.OwnTrack) []music.TrackSelect {
	return []music.TrackSelect{{
		trackList.Name,
		trackList.Album.Artist,
		trackList.Album.Album,
	},
	}
}
