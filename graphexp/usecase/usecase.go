package usecase

import (
	"math/rand"
	"musicAPI/graph/model"
	"musicAPI/graphexp"
	"strconv"
)

type useCase struct {
	PostgresRepo graphexp.Repository
}

func New(postgresRepo graphexp.Repository) graphexp.UseCase {
	return &useCase{PostgresRepo: postgresRepo}
}

func (uc *useCase) GetTracks() ([]*model.Tracks, error) {
	trackList, err := uc.PostgresRepo.GetTracks()
	if err != nil {
		return nil, err
	}
	return trackList, nil
}

func (uc *useCase) GetArtist(artist *string) ([]*model.Tracks, error) {
	trackList, err := uc.PostgresRepo.GetArtist(artist)
	if err != nil {
		return nil, err
	}
	return trackList, nil
}

func (uc *useCase) SetTrack(track model.NewTrack) error {
	list := rand.Intn(999)
	play := rand.Intn(999)
	newTrack := graphexp.OwnTrack{
		Name:      track.Name,
		Album:     track.Album,
		Artist:    track.Artist,
		Listeners: strconv.Itoa(list),
		Playcount: strconv.Itoa(play),
		Tag:       "Test",
	}
	err := uc.PostgresRepo.SetTrack(newTrack)
	if err != nil {
		return err
	}
	return nil
}
