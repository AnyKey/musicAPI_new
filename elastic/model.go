package elastic

import (
	"github.com/gorilla/websocket"
	"musicAPI/model"
)

type Repository interface {
	FullTextSearch(resData SocketSend) ([]model.TrackSelect, error)
}
type SocketSend struct {
	Track       string `json:"track"`
	NameCheck   bool   `json:"nameCheck"`
	ArtistCheck bool   `json:"artistCheck"`
	AlbumCheck  bool   `json:"albumCheck"`
}

type UseCase interface {
	WsSending(conn *websocket.Conn)
}
