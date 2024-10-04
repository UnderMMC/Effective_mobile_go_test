package repository

import (
	"EffectiveMobile_Go/internal/domain/entity"
	"database/sql"
	"strings"
)

type PostgresSongRepository struct {
	db *sql.DB
}

func NewPostgresSongRepository(db *sql.DB) *PostgresSongRepository {
	return &PostgresSongRepository{db: db}
}

func (r *PostgresSongRepository) GetAll(filter string) ([]entity.Song, error) {
	var songs []entity.Song
	var rows *sql.Rows
	var err error

	if filter == "" {
		query := `SELECT song_id, performer, song_name, release_date, lyric, link FROM music`
		rows, err = r.db.Query(query)
	} else {
		query := `SELECT song_id, performer, song_name, release_date, lyric, link FROM music WHERE CONCAT_WS(' ', song_id::text, performer, song_name, release_date, lyric, link) LIKE '%' || $1 || '%'`
		rows, err = r.db.Query(query, filter)
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var song entity.Song
		var songDetails entity.SongDetails
		err := rows.Scan(&song.ID, &song.Group, &song.Song, &songDetails.ReleaseDate, &songDetails.Text, &songDetails.Link)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func (r *PostgresSongRepository) Add(song entity.Song) error {
	emptyStr := ""
	query := `INSERT INTO music (performer, song_name, release_date, lyric, link) VALUES ($1, $2, $3, $4, $5)`
	err := r.db.QueryRow(query, song.Group, song.Song, emptyStr, emptyStr, emptyStr)
	if err != nil {
		return nil
	}
	return nil
}

func (r *PostgresSongRepository) Delete(group string, song string, id int) error {
	query := `DELETE FROM music WHERE performer=$1 AND song_name=$2 AND song_id=$3`
	_, err := r.db.Exec(query, group, song, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresSongRepository) Update(song entity.SongDetails, id int) error {
	query := `UPDATE music SET release_date = $1, lyric = $2, link = $3 WHERE song_id = $4`
	_, err := r.db.Exec(query, song.ReleaseDate, song.Text, song.Link, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresSongRepository) GetText(id, page, size int) ([]string, error) {
	var fullLyric string
	query := `SELECT lyric FROM music WHERE song_id=$1`
	err := r.db.QueryRow(query, id).Scan(&fullLyric)
	if err != nil {
		return nil, err
	}

	verses := strings.Split(fullLyric, "\n")

	start := page * size
	end := start + size

	if start >= len(verses) {
		return []string{}, nil
	}
	if end > len(verses) {
		end = len(verses)
	}

	return verses[start:end], nil
}

func (r *PostgresSongRepository) GetAllDetails(group, song string) (entity.SongDetails, error) {
	var songDetails entity.SongDetails
	query := `SELECT release_date, lyric, link FROM music WHERE song_name=$1 AND performer=$2`
	err := r.db.QueryRow(query, song, group).Scan(&songDetails.ReleaseDate, &songDetails.Text, &songDetails.Link)
	if err != nil {
		return songDetails, err
	}
	return songDetails, nil
}
