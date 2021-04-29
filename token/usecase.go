package token

import (
	"context"
	"musicAPI/model"
)

type UseCase interface {
	CheckToken(context.Context, string) bool
	NewToken(context.Context, string) (*model.Tokens, error)
	RefreshToken(context.Context, string) (*model.Tokens, error)
}
