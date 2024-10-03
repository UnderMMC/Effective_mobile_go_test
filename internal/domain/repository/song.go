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
		rows, err = r.db.Query("SELECT * FROM music")
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
		err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
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

// func (r *PostgresSongRepository) GetByID(id int) (entity.Song, error) {}

// func (r *PostgresSongRepository) Add(song entity.Song) (int, error) {}

// func (r *PostgresSongRepository) Update(song entity.Song) error {}

// func (r *PostgresSongRepository) Delete(id int) error {}
