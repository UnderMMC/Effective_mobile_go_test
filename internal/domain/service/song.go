package service

import "EffectiveMobile_Go/internal/domain/entity"

type SongRepository interface {
	GetAll(filter string) ([]entity.Song, error)
	Add(song entity.Song) (int, error)
	Delete(group string, song string, id int) error
	Update(song entity.SongDetails, id int) error
	//GetSongIGetSongID(song entity.Song) (int, error)
	// GetByID(id int) (entity.Song, error)
}

type SongService struct {
	songRepo SongRepository
}

func NewSongService(songRepo SongRepository) *SongService {
	return &SongService{songRepo: songRepo}
}

func (s *SongService) GetSongs(filter string) ([]entity.Song, error) {
	songs, err := s.songRepo.GetAll(filter)
	if err != nil {
		return nil, err
	}
	return songs, nil
}

func (s *SongService) AddSong(song entity.Song) (int, error) {
	var err error
	song.ID, err = s.songRepo.Add(song)
	if err != nil {
		return 0, err
	}
	return song.ID, nil
}

func (s *SongService) DeleteSong(group string, song string, id int) error {
	err := s.songRepo.Delete(group, song, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *SongService) UpdateSong(song entity.SongDetails, id int) error {
	err := s.songRepo.Update(song, id)
	if err != nil {
		return err
	}
	return nil
}
