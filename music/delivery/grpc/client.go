package grpc

import (
	"context"
	"encoding/json"
	pb "github.com/AnyKey/userslike/grpcsrv/like"
	"musicAPI/music"
	"time"
)

type Delivery struct {
	Conn pb.SubSrvClient
}

func New(conn pb.SubSrvClient) *Delivery {
	return &Delivery{
		Conn: conn,
	}
}

func (d *Delivery) SetLike(name, artist, token string) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	r, err := d.Conn.SetLike(ctx, &pb.LikeRequest{
		Name:   name,
		Artist: artist,
		Jwt:    token,
	})
	if err != nil {
		return nil, err
	}
	defer cancel()
	res := r.GetMessage()
	return &res, nil
}

func (d *Delivery) GetLike(name, artist, token string) (*music.LikeList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	r, err := d.Conn.GetLike(ctx, &pb.TrackRequest{
		Name:   name,
		Artist: artist,
		Jwt:    token,
	})
	if err != nil {
		return nil, err
	}
	var likes []music.LikeSelect
	err = json.Unmarshal(r.GetUser(), &likes)
	if err != nil {
		return nil, err
	}
	likeList := music.LikeList{
		Track:  r.GetName(),
		Artist: r.GetArtist(),
		Likes:  likes,
	}
	defer cancel()
	return &likeList, nil
}
