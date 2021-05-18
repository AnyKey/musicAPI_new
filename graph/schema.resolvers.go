package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"musicAPI/graph/generated"
	"musicAPI/graph/model"
)

func (r *mutationResolver) CreateTrack(ctx context.Context, input model.NewTrack) (*model.Tracks, error) {
	result := model.Tracks{Album: input.Album, Artist: input.Artist, Name: input.Name}
	r.ExpUseCase.SetTrack(input)
	return &result, nil
}

func (r *queryResolver) Tracks(ctx context.Context) ([]*model.Tracks, error) {
	result, err := r.ExpUseCase.GetTracks()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *queryResolver) Artist(ctx context.Context, artist *string) ([]*model.Tracks, error) {
	result, err := r.ExpUseCase.GetArtist(artist)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
