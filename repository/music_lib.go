package repository

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"log"
	"musicAPI/handlers/request"
	"musicAPI/model"
)

type Repository struct {
	Conn *sql.DB
}

func (repo Repository) GetTracks(track string, artist string) ([]model.TrackSelect, error) {
	var trackList []model.TrackSelect
	rows, err := repo.Conn.Query("SELECT name, artist, listeners FROM tracklist WHERE name = %$1% and artist = $2", track, artist)
	if err != nil {
		return nil, errors.Wrap(err, "error select in DB!")
	}
	defer rows.Close()

	for rows.Next() {
		tl := model.TrackSelect{}
		err := rows.Scan(&tl.Name, &tl.Artist, &tl.Listeners)
		if err != nil {
			return nil, errors.Wrap(err, "error Scan values")
		}
		trackList = append(trackList, tl)
	}
	return trackList, nil
}

func (repo Repository) SetTracks(NewTracks request.OwnTrack) error {
	ctx := context.Background()
	tx, err := repo.Conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := tx.ExecContext(ctx, "INSERT INTO album (name) VALUES ($1)", NewTracks.Album.Album)
	if err != nil {
		tx.Rollback()
		fmt.Println("ALBUM!", err.Error())
		return err
	}
	albumId, _ := res.LastInsertId()

	res, err = tx.ExecContext(ctx, "INSERT INTO artist (name) VALUES ($1)", NewTracks.Album.Artist)
	if err != nil {
		fmt.Println("ARTIST", err.Error())
		tx.Rollback()
		return err
	}
	artistId, _ := res.LastInsertId()
	fmt.Println(artistId)
	_, err = tx.ExecContext(ctx, "INSERT INTO tag (genre) VALUES ($1)", NewTracks.TopTags.Genre[0].Tag)
	if err != nil {
		fmt.Println("TAG!", err.Error())
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO track (name, artist_id, album_id, listeners, playcount, tag) VALUES ($1, $2, $3, $4, $5, $6)", NewTracks.Name, artistId, albumId, NewTracks.Listeners, NewTracks.Playcount, NewTracks.TopTags.Genre[0].Tag)
	if err != nil {
		fmt.Println("TRACK!", err.Error())
		tx.Rollback()
		return err
	}
	//QueryRows прочекать или return (в синтаксисе postgres)
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	//for i:= range NewTracks{
	//		listeners, _ := strconv.Atoi(NewTracks[i].Listeners)
	//		_, err := repo.Conn.Exec("INSERT INTO tracklist (name, artist, listeners) VALUES ($1, $2, $3)", NewTracks[i].Name, NewTracks[i].Artist, listeners)
	//		if err != nil {
	//			return errors.Wrap(err, "error insert in DB!")
	//		}
	//	}
	return nil
}

func (repo Repository) GetArtistTracks(artist string) ([]model.TrackSelect, error) {
	var trackList []model.TrackSelect
	rows, err := repo.Conn.Query("SELECT name, artist, listeners FROM tracklist WHERE artist = $1", artist)
	if err != nil {
		return nil, errors.Wrap(err, "error select in DB!")
	}
	defer rows.Close()

	for rows.Next() {
		tl := model.TrackSelect{}
		err := rows.Scan(&tl.Name, &tl.Artist, &tl.Listeners)
		if err != nil {
			return nil, errors.Wrap(err, "error Scan values")
		}
		trackList = append(trackList, tl)
	}
	return trackList, nil
}
