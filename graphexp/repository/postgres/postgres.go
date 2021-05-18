package postgres

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"log"
	"musicAPI/graph/model"
	"musicAPI/graphexp"
)

type Repository struct {
	Conn *sql.DB
}

func New(conn *sql.DB) *Repository {
	return &Repository{
		Conn: conn,
	}
}

func (repo *Repository) GetTracks() ([]*model.Tracks, error) {

	var trackList []*model.Tracks
	rows, err := repo.Conn.Query("SELECT track.name as track, artist.name as artist, album.name as album  FROM track, artist, album " +
		"WHERE track.artist_id = artist.id and track.album_id = album.id")
	if err != nil {
		return nil, errors.Wrap(err, "error select in DB!")
	}
	defer rows.Close()

	for rows.Next() {
		tl := model.Tracks{}
		err := rows.Scan(&tl.Name, &tl.Artist, &tl.Album)
		if err != nil {
			return nil, errors.Wrap(err, "error Scan values")
		}
		trackList = append(trackList, &tl)
	}

	return trackList, nil
}

func (repo *Repository) GetArtist(artist *string) ([]*model.Tracks, error) {

	var trackList []*model.Tracks
	rows, err := repo.Conn.Query("SELECT track.name as track, artist.name as artist, album.name as album  FROM track, artist, album "+
		"WHERE track.artist_id = artist.id and track.album_id = album.id and artist.name = $1", artist)
	if err != nil {
		return nil, errors.Wrap(err, "error select in DB!")
	}
	defer rows.Close()

	for rows.Next() {
		tl := model.Tracks{}
		err := rows.Scan(&tl.Name, &tl.Artist, &tl.Album)
		if err != nil {
			return nil, errors.Wrap(err, "error Scan values")
		}
		trackList = append(trackList, &tl)
	}

	return trackList, nil
}
func (repo *Repository) SetTrack(newTracks graphexp.OwnTrack) error {
	ctx := context.Background()
	tx, err := repo.Conn.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
	}
	// defer commit rollback tnx
	var lastinsertedid int
	rows, err := tx.QueryContext(ctx, "SELECT id FROM album WHERE name = $1", newTracks.Album)
	if err != nil {
		return errors.Wrap(err, "error select in DB!")
	}
	var albumId int
	if rows != nil {
		var tl int
		for rows.Next() {
			err = rows.Scan(&tl)
			if err != nil {
				return errors.Wrap(err, "error Scan values")
			}
		}
		albumId = tl
	}
	if albumId == 0 {
		err = tx.QueryRowContext(ctx, "INSERT INTO album (name) VALUES ($1) returning id", newTracks.Album).Scan(&lastinsertedid)
		if err != nil {
			tx.Rollback()
			log.Println("ALBUM!", err.Error())
			return err
		}
		albumId = lastinsertedid
	}
	rows.Close()
	if albumId == 0 {
		tx.Rollback()
		return errors.New("Empty album")
	}
	rows, err = tx.QueryContext(ctx, "SELECT id FROM artist WHERE name = $1", newTracks.Artist)
	if err != nil {
		return errors.Wrap(err, "error select in DB!")
	}
	var artistId int
	if rows != nil {
		var tl int
		for rows.Next() {
			err = rows.Scan(&tl)
			if err != nil {
				return errors.Wrap(err, "error Scan values")
			}
		}
		artistId = tl
	}
	if artistId == 0 {
		err = tx.QueryRowContext(ctx, "INSERT INTO artist (name) VALUES ($1) returning id", newTracks.Artist).Scan(&lastinsertedid)
		if err != nil {
			log.Println("ARTIST", err.Error())
			err = tx.Rollback()
			if err != nil {
				log.Println(err)
			}
			return err
		}
		artistId = lastinsertedid
	}
	rows.Close()
	if artistId == 0 {
		tx.Rollback()
		return errors.New("Empty artist")
	}
	checkTag := newTracks.Tag
	if len(checkTag) == 0 {
		return errors.New("Empty tags")
	}

	rows, err = tx.QueryContext(ctx, "SELECT genre FROM tag WHERE genre = $1", newTracks.Tag)
	if err != nil {
		return errors.Wrap(err, "error select in DB!")
	}
	var tag string
	if rows != nil {
		for rows.Next() {
			err = rows.Scan(&tag)
			if err != nil {
				return errors.Wrap(err, "error Scan values")
			}
		}
	}
	if tag == "" {
		_, err = tx.ExecContext(ctx, "INSERT INTO tag (genre) VALUES ($1)", newTracks.Tag)
		if err != nil {
			log.Println("TAG!", err.Error())
			tx.Rollback()
			return err
		}
	}
	rows.Close()

	_, err = tx.ExecContext(ctx, "INSERT INTO track (name, artist_id, album_id, listeners, playcount, tag) VALUES ($1, $2, $3, $4, $5, $6)", newTracks.Name, artistId, albumId, newTracks.Listeners, newTracks.Playcount, newTracks.Tag)
	if err != nil && err == sql.ErrNoRows {
		log.Println("TRACK!", err.Error())
		tx.Rollback()
		return errors.Wrap(err, "transaction cancel")
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
	}

	return nil
}
