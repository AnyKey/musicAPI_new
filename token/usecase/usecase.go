package usecase

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"log"
	"musicAPI/model"
	"musicAPI/token"
	"time"
)

type tokenUseCase struct { // tokenUsecase {}
	TokenRepo token.TokenRepository
}

/* func NewTokenUseCase(tokenRepo token.TokenRepository) *TokenUseCase {
	return &TokenUseCase{TokenRepo: tokenRepo}
} */

func New(tokenRepo token.TokenRepository) token.UseCase {
	return &tokenUseCase{TokenRepo: tokenRepo}
}

//	CheckToken(ctx context.Context, token string) bool
//	NewToken(ctx context.Context, user string) (*model.Tokens, error)
//	RefreshToken(ctx context.Context, token string) (*model.Tokens, error)

func (t *tokenUseCase) CheckToken(ctx context.Context, myToken string) bool {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(myToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("key"), nil
	})

	if err == nil && token.Valid {
		user := (claims["name"]).(string)
		res := t.TokenRepo.GetToken(ctx, user)
		if res.Access == myToken && res.Valid == true {
			return true
		}
	}
	return false
}

func (t *tokenUseCase) NewToken(ctx context.Context, user string) (*model.Tokens, error) {
	tokenA, err := newTokenA(user)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	tokenR, err := newTokenR(user)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	access, err := tokenA.SignedString([]byte("key"))
	refresh, err := tokenR.SignedString([]byte("key"))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	newTokens := model.Tokens{
		Access:  access,
		Refresh: refresh,
		Valid:   true,
	}
	err = t.TokenRepo.SetToken(ctx, user, newTokens)
	return &newTokens, nil
}

func (t *tokenUseCase) RefreshToken(ctx context.Context, rToken string) (*model.Tokens, error) {

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(rToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("key"), nil
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if !token.Valid {
		return nil, errors.Wrap(err, "invalid refresh token")
	}
	user := claims["name"].(string)
	res := t.TokenRepo.GetToken(ctx, user)

	if res.Refresh != rToken {
		return nil, errors.Wrap(err, "invalid refresh token")
	}
	tokenA, err := newTokenA(user)
	if err != nil {
		return nil, err
	}
	tokenR, err := newTokenR(user)
	if err != nil {
		return nil, err
	}
	access, err := tokenA.SignedString([]byte("key"))
	refresh, err := tokenR.SignedString([]byte("key"))
	newToken := model.Tokens{
		Access:  access,
		Refresh: refresh,
		Valid:   true,
	}
	if err != nil {
		return nil, err
	}
	err = t.TokenRepo.SetToken(ctx, user, newToken)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &newToken, nil
}

func newTokenA(user string) (*jwt.Token, error) {
	return jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		jwt.MapClaims{
			"name": user,
			"exp":  time.Now().Add(time.Hour * 3).Unix(),
			"root": true,
		}), nil
}
func newTokenR(user string) (*jwt.Token, error) {
	return jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		jwt.MapClaims{
			"name": user,
			"exp":  time.Now().Add(time.Hour * 6).Unix(),
		}), nil
}
