package usecase

import (
	"musicAPI/graph/model"
	"musicAPI/graphexp"
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
