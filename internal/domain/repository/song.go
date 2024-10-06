package repository

import (
	"EffectiveMobile_Go/internal/domain/entity"
	"database/sql"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"strings"
)

type PostgresSongRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPostgresSongRepository(db *sql.DB, logger *zap.Logger) *PostgresSongRepository {
	return &PostgresSongRepository{db: db, logger: logger}
}

func (r *PostgresSongRepository) GetAll(filter string) ([]entity.Song, error) {
	r.logger.Debug("Fetching all songs", zap.String("filter", filter))

	var songs []entity.Song
	var rows *sql.Rows
	var err error

	if filter == "" {
		query := `SELECT song_id, performer, song_name, release_date, lyric, link FROM music`
		rows, err = r.db.Query(query)
	} else {
		query := `SELECT song_id, performer, song_name, release_date, lyric, link
				  FROM music 
				  WHERE CONCAT_WS(' ', song_id::text, performer, song_name, release_date, lyric, link) 
				  LIKE '%' || $1 || '%'`
		rows, err = r.db.Query(query, filter)
	}
	if err != nil {
		r.logger.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var song entity.Song
		var songDetails entity.SongDetails
		err := rows.Scan(&song.ID, &song.Group, &song.Song, &songDetails.ReleaseDate, &songDetails.Text, &songDetails.Link)
		if err != nil {
			r.logger.Error("Failed to scan row", zap.Error(err))
			return nil, err
		}
		songs = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("Error encountered while reading rows", zap.Error(err))
		return nil, err
	}

	r.logger.Info("Fetched songs successfully", zap.Int("count", len(songs)))
	return songs, nil
}

func (r *PostgresSongRepository) Add(song entity.Song) error {
	r.logger.Debug("Adding new song", zap.Any("song", song))

	emptyStr := ""
	query := `INSERT INTO music (performer, song_name, release_date, lyric, link) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, song.Group, song.Song, emptyStr, emptyStr, emptyStr)
	if err != nil {
		r.logger.Error("Failed to add song", zap.Error(err))
		return err
	}
	r.logger.Info("Song added successfully", zap.String("song", song.Song), zap.String("group", song.Group))
	return nil
}

func (r *PostgresSongRepository) Delete(group string, song string, id int) error {
	r.logger.Debug("Deleting song", zap.String("group", group), zap.String("song", song), zap.Int("id", id))

	query := `DELETE FROM music WHERE performer=$1 AND song_name=$2 AND song_id=$3`
	_, err := r.db.Exec(query, group, song, id)
	if err != nil {
		r.logger.Error("Failed to delete song", zap.Error(err))
		return err
	}
	r.logger.Info("Song deleted successfully", zap.String("group", group), zap.String("song", song))
	return nil
}

func (r *PostgresSongRepository) Update(song entity.SongDetails, id int) error {
	r.logger.Debug("Updating song", zap.Any("songDetails", song), zap.Int("id", id))

	query := `UPDATE music SET release_date = $1, lyric = $2, link = $3 WHERE song_id = $4`
	_, err := r.db.Exec(query, song.ReleaseDate, song.Text, song.Link, id)
	if err != nil {
		r.logger.Error("Failed to update song", zap.Error(err))
		return err
	}
	r.logger.Info("Song updated successfully", zap.Int("id", id))
	return nil
}

func (r *PostgresSongRepository) GetText(id, page, size int) ([]string, error) {
	r.logger.Debug("Fetching song lyrics", zap.Int("id", id), zap.Int("page", page), zap.Int("size", size))

	var fullLyric string
	query := `SELECT lyric FROM music WHERE song_id=$1`
	err := r.db.QueryRow(query, id).Scan(&fullLyric)
	if err != nil {
		r.logger.Error("Failed to fetch song lyrics", zap.Error(err))
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

	r.logger.Info("Fetched lyrics successfully", zap.Int("count", len(verses[start:end])))
	return verses[start:end], nil
}

func (r *PostgresSongRepository) GetAllDetails(group, song string) (entity.SongDetails, error) {
	r.logger.Debug("Fetching all details for song", zap.String("group", group), zap.String("song", song))

	var songDetails entity.SongDetails
	query := `SELECT release_date, lyric, link FROM music WHERE song_name=$1 AND performer=$2`
	err := r.db.QueryRow(query, song, group).Scan(&songDetails.ReleaseDate, &songDetails.Text, &songDetails.Link)
	if err != nil {
		r.logger.Error("Failed to fetch song details", zap.Error(err))
		return songDetails, err
	}
	r.logger.Info("Fetched song details successfully", zap.Any("songDetails", songDetails))
	return songDetails, nil
}
