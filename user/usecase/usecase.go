package usecase

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"log"
	"musicAPI/user"
	"time"
)

type tokenUseCase struct {
	TokenRepo user.Repository
}

func New(tokenRepo user.Repository) user.UseCase {
	return &tokenUseCase{TokenRepo: tokenRepo}
}

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

func (t *tokenUseCase) NewToken(ctx context.Context, username string) (*user.Tokens, error) {
	tokenA, err := newTokenAccess(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	tokenR, err := newTokenRefresh(username)
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
	newTokens := user.Tokens{
		Access:  access,
		Refresh: refresh,
		Valid:   true,
	}
	err = t.TokenRepo.SetToken(ctx, username, newTokens)
	return &newTokens, nil
}

func (t *tokenUseCase) RefreshToken(ctx context.Context, refreshToken string) (*user.Tokens, error) {

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("key"), nil
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if !token.Valid {
		return nil, errors.Wrap(err, "invalid refresh token")
	}
	username := claims["name"].(string)
	res := t.TokenRepo.GetToken(ctx, username)

	if res.Refresh != refreshToken {
		return nil, errors.Wrap(err, "invalid refresh token")
	}
	tokenAccess, err := newTokenAccess(username)
	if err != nil {
		return nil, err
	}
	tokenRefresh, err := newTokenRefresh(username)
	if err != nil {
		return nil, err
	}
	access, err := tokenAccess.SignedString([]byte("key"))
	refresh, err := tokenRefresh.SignedString([]byte("key"))
	newToken := user.Tokens{
		Access:  access,
		Refresh: refresh,
		Valid:   true,
	}
	if err != nil {
		return nil, err
	}
	err = t.TokenRepo.SetToken(ctx, username, newToken)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &newToken, nil
}

func newTokenAccess(user string) (*jwt.Token, error) {
	return jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		jwt.MapClaims{
			"name": user,
			"exp":  time.Now().Add(time.Hour * 3).Unix(),
			"root": true,
		}), nil
}
func newTokenRefresh(user string) (*jwt.Token, error) {
	return jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		jwt.MapClaims{
			"name": user,
			"exp":  time.Now().Add(time.Hour * 6).Unix(),
		}), nil
}
