package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"musicAPI/music"
)

type Repository struct {
	Conn *sql.DB
}

func New(conn *sql.DB) *Repository {
	return &Repository{
		Conn: conn,
	}
}

func (repo *Repository) GetTracks(track string, artist string) ([]music.TrackSelect, error) {

	var trackList []music.TrackSelect
	rows, err := repo.Conn.Query("SELECT track.name as track, artist.name as artist, album.name as album  FROM track, artist, album "+
		"WHERE track.artist_id = artist.id and track.album_id = album.id and track.name = $1 AND artist.name = $2", track, artist)
	if err != nil {
		return nil, errors.Wrap(err, "error select in DB!")
	}
	defer rows.Close()

	for rows.Next() {
		tl := music.TrackSelect{}
		err := rows.Scan(&tl.Name, &tl.Artist, &tl.Album)
		if err != nil {
			return nil, errors.Wrap(err, "error Scan values")
		}
		trackList = append(trackList, tl)
	}

	return trackList, nil
}

func (repo *Repository) CheckTrack(track string, artist string) (bool, error) {

	var trackList []music.TrackSelect
	rows, err := repo.Conn.Query("SELECT track.name as track, artist.name as artist, album.name as album  FROM track, artist, album "+
		"WHERE track.artist_id = artist.id and track.album_id = album.id and track.name = $1 AND artist.name = $2", track, artist)
	if err != nil {
		return false, errors.Wrap(err, "error select in DB!")
	}
	defer rows.Close()

	for rows.Next() {
		tl := music.TrackSelect{}
		err := rows.Scan(&tl.Name, &tl.Artist, &tl.Album)
		if err != nil {
			return false, errors.Wrap(err, "error Scan values")
		}
		trackList = append(trackList, tl)
	}
	if trackList == nil {
		return false, nil
	}
	return true, nil
}

func (repo *Repository) SetTracks(newTracks music.OwnTrack) error {
	ctx := context.Background()
	tx, err := repo.Conn.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
	}
	// defer commit rollback tnx
	var lastinsertedid int
	rows, err := tx.QueryContext(ctx, "SELECT id FROM album WHERE name = $1", newTracks.Album.Album)
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
		err = tx.QueryRowContext(ctx, "INSERT INTO album (name) VALUES ($1) returning id", newTracks.Album.Album).Scan(&lastinsertedid)
		if err != nil {
			tx.Rollback()
			fmt.Println("ALBUM!", err.Error())
			return err
		}
		albumId = lastinsertedid
	}
	rows.Close()
	if albumId == 0 {
		tx.Rollback()
		return errors.New("Empty album")
	}
	rows, err = tx.QueryContext(ctx, "SELECT id FROM artist WHERE name = $1", newTracks.Album.Artist)
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
		err = tx.QueryRowContext(ctx, "INSERT INTO artist (name) VALUES ($1) returning id", newTracks.Album.Artist).Scan(&lastinsertedid)
		if err != nil {
			fmt.Println("ARTIST", err.Error())
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
	checkTag := newTracks.TopTags.Genre
	if len(checkTag) == 0 {
		return errors.New("Empty tags")
	}

	rows, err = tx.QueryContext(ctx, "SELECT genre FROM tag WHERE genre = $1", newTracks.TopTags.Genre[0].Tag)
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
		_, err = tx.ExecContext(ctx, "INSERT INTO tag (genre) VALUES ($1)", newTracks.TopTags.Genre[0].Tag)
		if err != nil {
			fmt.Println("TAG!", err.Error())
			tx.Rollback()
			return err
		}
	}
	rows.Close()

	_, err = tx.ExecContext(ctx, "INSERT INTO track (name, artist_id, album_id, listeners, playcount, tag) VALUES ($1, $2, $3, $4, $5, $6)", newTracks.Name, artistId, albumId, newTracks.Listeners, newTracks.Playcount, newTracks.TopTags.Genre[0].Tag)
	if err != nil && err == sql.ErrNoRows {
		fmt.Println("TRACK!", err.Error())
		tx.Rollback()
		return nil
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (repo *Repository) GetGenreTracks(genre string) ([]music.TrackSelect, error) {

	var trackList []music.TrackSelect
	rows, err := repo.Conn.Query("SELECT track.name as track, artist.name as artist, album.name as album "+
		"FROM track, artist, album "+
		"WHERE track.artist_id = artist.id and track.album_id = album.id and track.tag = $1", genre)
	if err != nil {
		return nil, errors.Wrap(err, "error select in DB!")
	}
	defer rows.Close()

	for rows.Next() {
		tl := music.TrackSelect{}
		err := rows.Scan(&tl.Name, &tl.Artist, &tl.Album)
		if err != nil {
			return nil, errors.Wrap(err, "error Scan values")
		}
		trackList = append(trackList, tl)
	}
	if err != nil || trackList == nil {
		return nil, errors.New("empty values")
	}

	return trackList, nil
}

func (repo *Repository) GetArtistTracks(artist string) ([]music.TrackSelect, error) {

	var trackList []music.TrackSelect
	rows, err := repo.Conn.Query("SELECT track.name as track, artist.name as artist, album.name as album  FROM track, artist, album "+
		"WHERE track.artist_id = artist.id and track.album_id = album.id and artist.name = $1", artist)

	if err != nil {
		return nil, errors.Wrap(err, "error select in DB!")
	}
	defer rows.Close()

	for rows.Next() {
		tl := music.TrackSelect{}
		err := rows.Scan(&tl.Name, &tl.Artist, &tl.Album)
		if err != nil {
			return nil, errors.Wrap(err, "error Scan values")
		}
		trackList = append(trackList, tl)
	}
	if err != nil || trackList == nil {
		return nil, errors.New("empty values")
	}
	return trackList, nil
}

func (repo *Repository) GetChart(sortTo string) ([]music.ChartSelect, error) {

	var trackList []music.ChartSelect
	var querySql string
	if sortTo == "list" {
		querySql = "SELECT track.listeners as listeners, track.playcount as playcount, track.name as track," +
			" artist.name as artist, album.name as album,  track.tag as genre FROM track, artist, album " +
			"WHERE track.artist_id = artist.id and track.album_id = album.id ORDER BY listeners desc "
	} else if sortTo == "play" {
		querySql = "SELECT track.listeners as listeners, track.playcount as playcount, track.name as track," +
			" artist.name as artist, album.name as album,  track.tag as genre FROM track, artist, album " +
			"WHERE track.artist_id = artist.id and track.album_id = album.id ORDER BY playcount desc "
	} else {
		fmt.Println("there is not such method -->", sortTo)
		return nil, nil
	}

	rows, err := repo.Conn.Query(querySql)
	if err != nil {
		return nil, errors.Wrap(err, "error select in DB!")
	}
	defer rows.Close()

	for rows.Next() {
		tl := music.ChartSelect{}
		err := rows.Scan(&tl.Listeners, &tl.Playcount, &tl.Track, &tl.Artist, &tl.Album, &tl.Genre)
		if err != nil {
			return nil, errors.Wrap(err, "error Scan values")
		}
		trackList = append(trackList, tl)
	}

	return trackList, nil
}
