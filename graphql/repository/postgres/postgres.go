package postgres

import (
	"database/sql"
	"github.com/pkg/errors"
	"musicAPI/graphql"
)

type Repository struct {
	Conn *sql.DB
}

func New(conn *sql.DB) *Repository {
	return &Repository{
		Conn: conn,
	}
}

func (repo *Repository) GetTracks() ([]graphql.TrackSelect, error) {

	var trackList []graphql.TrackSelect
	rows, err := repo.Conn.Query("SELECT track.name as track, artist.name as artist, album.name as album  FROM track, artist, album " +
		"WHERE track.artist_id = artist.id and track.album_id = album.id")
	if err != nil {
		return nil, errors.Wrap(err, "error select in DB!")
	}
	defer rows.Close()

	for rows.Next() {
		tl := graphql.TrackSelect{}
		err := rows.Scan(&tl.Name, &tl.Artist, &tl.Album)
		if err != nil {
			return nil, errors.Wrap(err, "error Scan values")
		}
		trackList = append(trackList, tl)
	}

	return trackList, nil
}
