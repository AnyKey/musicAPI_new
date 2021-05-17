package usecase

import (
	"musicAPI/graphql"
)

type useCase struct {
	GraphRepo graphql.Repository
}

func New(graphRepo graphql.Repository) graphql.UseCase {
	return &useCase{GraphRepo: graphRepo}
}

func (uc *useCase) GetTracks() ([]graphql.TrackSelect, error) {
	trackList, err := uc.GraphRepo.GetTracks()
	if err != nil {
		return nil, err
	}

	return trackList, nil
}
