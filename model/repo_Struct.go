package model

type TrackSelect struct {
	Name     string `db:"name" json:"name"`
	Artist   string `db:"artist" json:"artist"`
	Listeners int `db:"listeners" json:"listeners"`
}

type TrackAppend struct {

}
