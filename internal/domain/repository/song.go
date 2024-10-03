package repository

import (
	"EffectiveMobile_Go/internal/domain/entity"
	"database/sql"
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
		rows, err = r.db.Query("SELECT * FROM music") // ORDER BY song_id LIMIT $1 OFFSET $2", limit, offset)
	} else {
		query := `SELECT * FROM music WHERE CONCAT_WS(' ', song_id::text, performer, song_name, release_data, lyric, link) LIKE '%' || $1 || '%'`
		rows, err = r.db.Query(query, filter)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song entity.Song
		var songd entity.SongDetails
		err := rows.Scan(&song.ID, &song.Group, &song.Song, &songd.ReleaseDate, &songd.Text, &songd.Link)
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

func (r *PostgresSongRepository) Add(song entity.Song) (int, error) {
	err := r.db.QueryRow("INSERT INTO music (performer, song_name) VALUES ($1, $2) RETURNING song_id", song.Group, song.Song).Scan(&song.ID)
	if err != nil {
		return 0, err
	}
	return song.ID, nil
}

//func (r *PostgresSongRepository) GetSongID(song entity.Song) (int, error) {
//	err := r.db.QueryRow("SELECT id FROM music WHERE performer=$1 AND song_name=$2", song.Group, song.Song).Scan(&song.ID)
//	if err != nil {
//		return 0, err
//	}
//	return song.ID, err
//}

// func (r *PostgresSongRepository) GetByID(id int) (entity.Song, error) {}

// func (r *PostgresSongRepository) Update(song entity.Song) error {}

// func (r *PostgresSongRepository) Delete(id int) error {}
