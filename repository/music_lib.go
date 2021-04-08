package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"log"
	"musicAPI/model"
	"time"
)

type Repository struct {
	Conn  *sql.DB
	Redis *redis.Client
}

func (repo Repository) GetTracks(track string, artist string) ([]model.TrackSelect, error) {
	var ctx = context.Background()
	result := repo.GetTracksRedis(track, artist)
	if result != nil {
		return result, nil
	}
	var trackList []model.TrackSelect
	rows, err := repo.Conn.Query("SELECT track.name as track, artist.name as artist, album.name as album  FROM track, artist, album "+
		"WHERE track.artist_id = artist.id and track.album_id = album.id and track.name = $1 AND artist.name = $2", track, artist)
	if err != nil {
		return nil, errors.Wrap(err, "error select in DB!")
	}
	defer rows.Close()

	for rows.Next() {
		tl := model.TrackSelect{}
		err := rows.Scan(&tl.Name, &tl.Artist, &tl.Album)
		if err != nil {
			return nil, errors.Wrap(err, "error Scan values")
		}
		trackList = append(trackList, tl)
	}
	bytes, err := json.Marshal(trackList)
	if err == nil {
		repo.Redis.Set(ctx, "Track:"+track+"_Artist:"+artist, bytes, 5*time.Minute)
	}
	return trackList, nil
}

func (repo Repository) SetTracks(NewTracks model.OwnTrack) error {
	ctx := context.Background()
	tx, err := repo.Conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	// defer commit rollback tnx
	var lastinsertedid int
	rows, err := tx.QueryContext(ctx, "SELECT id FROM album WHERE name = $1", NewTracks.Album.Album)
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
		err = tx.QueryRowContext(ctx, "INSERT INTO album (name) VALUES ($1) returning id", NewTracks.Album.Album).Scan(&lastinsertedid)
		if err != nil {
			tx.Rollback()
			fmt.Println("ALBUM!", err.Error())
			return err
		}
		albumId = lastinsertedid
	}
	rows.Close()

	rows, err = tx.QueryContext(ctx, "SELECT id FROM artist WHERE name = $1", NewTracks.Album.Artist)
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
		err = tx.QueryRowContext(ctx, "INSERT INTO artist (name) VALUES ($1) returning id", NewTracks.Album.Artist).Scan(&lastinsertedid)
		if err != nil {
			fmt.Println("ARTIST", err.Error())
			tx.Rollback()
			return err
		}
		artistId = lastinsertedid
	}
	rows.Close()

	rows, err = tx.QueryContext(ctx, "SELECT genre FROM tag WHERE genre = $1", NewTracks.TopTags.Genre[0].Tag)
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
		_, err = tx.ExecContext(ctx, "INSERT INTO tag (genre) VALUES ($1)", NewTracks.TopTags.Genre[0].Tag)
		if err != nil {
			fmt.Println("TAG!", err.Error())
			tx.Rollback()
			return err
		}
	}
	rows.Close()

	_, err = tx.ExecContext(ctx, "INSERT INTO track (name, artist_id, album_id, listeners, playcount, tag) VALUES ($1, $2, $3, $4, $5, $6)", NewTracks.Name, artistId, albumId, NewTracks.Listeners, NewTracks.Playcount, NewTracks.TopTags.Genre[0].Tag)
	if err != nil {
		fmt.Println("TRACK!", err.Error())
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (repo Repository) GetGenreTracks(genre string) ([]model.TrackSelect, error) {

	var ctx = context.Background()
	result := repo.GetGenreRedis(genre)
	if result != nil {
		return result, nil
	}
	var trackList []model.TrackSelect
	rows, err := repo.Conn.Query("SELECT track.name as track, artist.name as artist, album.name as album "+
		"FROM track, artist, album "+
		"WHERE track.artist_id = artist.id and track.album_id = album.id and track.tag = $1", genre)
	if err != nil {
		return nil, errors.Wrap(err, "error select in DB!")
	}
	defer rows.Close()

	for rows.Next() {
		tl := model.TrackSelect{}
		err := rows.Scan(&tl.Name, &tl.Artist, &tl.Album)
		if err != nil {
			return nil, errors.Wrap(err, "error Scan values")
		}
		trackList = append(trackList, tl)
	}
	bytes, err := json.Marshal(trackList)
	if err == nil {
		repo.Redis.Set(ctx, "Genre:"+genre, bytes, 5*time.Minute)
	}

	return trackList, nil
}

func (repo Repository) GetArtistTracks(artist string) ([]model.TrackSelect, error) {

	var ctx = context.Background()
	result := repo.GetArtistRedis(artist)
	if result != nil {
		return result, nil
	}
	var trackList []model.TrackSelect
	rows, err := repo.Conn.Query("SELECT track.name as track, artist.name as artist, album.name as album  FROM track, artist, album "+
		"WHERE track.artist_id = artist.id and track.album_id = album.id and artist.name = $1", artist)
	log.Println("err", err)

	if err != nil {
		return nil, errors.Wrap(err, "error select in DB!")
	}
	defer rows.Close()

	for rows.Next() {
		tl := model.TrackSelect{}
		err := rows.Scan(&tl.Name, &tl.Artist, &tl.Album)
		if err != nil {
			return nil, errors.Wrap(err, "error Scan values")
		}
		trackList = append(trackList, tl)
	}

	bytes, err := json.Marshal(trackList)
	if err == nil {
		repo.Redis.Set(ctx, "Artist:"+artist, bytes, 5*time.Minute)
	}
	return trackList, nil
}

func (repo Repository) GetChart(sortTo string) ([]model.ChartSelect, error) {

	var ctx = context.Background()
	result := repo.GetChartRedis(sortTo)
	if result != nil {
		return result, nil
	}
	var trackList []model.ChartSelect
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
		tl := model.ChartSelect{}
		err := rows.Scan(&tl.Listeners, &tl.Playcount, &tl.Track, &tl.Artist, &tl.Album, &tl.Genre)
		if err != nil {
			return nil, errors.Wrap(err, "error Scan values")
		}
		trackList = append(trackList, tl)
	}
	bytes, err := json.Marshal(trackList)
	if err == nil {
		repo.Redis.Set(ctx, "SortTo:"+sortTo, bytes, 5*time.Minute)
	}
	return trackList, nil
}
