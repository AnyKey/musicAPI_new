package token

import (
	"context"
	"musicAPI/model"
)

type TokenRepository interface {
	GetToken(ctx context.Context, user string) *model.Tokens
	SetToken(ctx context.Context, user string, tokens model.Tokens) error
}
