package elastic

import (
	"github.com/gorilla/websocket"
)

type Delivery interface {
	FullTextSearch(SocketSend) ([]TrackSelect, error)
}
type SocketSend struct {
	Track       string `json:"track"`
	NameCheck   bool   `json:"nameCheck"`
	ArtistCheck bool   `json:"artistCheck"`
	AlbumCheck  bool   `json:"albumCheck"`
}

type UseCase interface {
	WsSending(*websocket.Conn)
}

type TrackSelect struct {
	Name   string `db:"track" json:"name"`
	Artist string `db:"artist" json:"artist"`
	Album  string `db:"album" json:"album"`
}
