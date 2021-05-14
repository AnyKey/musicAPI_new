package user

import (
	"context"
)

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
	Valid   bool   `json:"valid"`
}
type UseCase interface {
	CheckToken(context.Context, string) bool
	NewToken(context.Context, string) (*Tokens, error)
	RefreshToken(context.Context, string) (*Tokens, error)
}
type Repository interface {
	GetToken(context.Context, string) *Tokens
	SetToken(context.Context, string, Tokens) error
}
