package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"musicAPI/graph/generated"
	"musicAPI/graph/model"
)

func (r *queryResolver) Tracks(ctx context.Context) ([]*model.Tracks, error) {
	result, err := r.ExpUseCase.GetTracks()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
