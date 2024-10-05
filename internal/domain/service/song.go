package service

import "EffectiveMobile_Go/internal/domain/entity"

type SongRepository interface {
	GetAll(filter string) ([]entity.Song, error)
	Add(song entity.Song) error
	Delete(group string, song string, id int) error
	Update(song entity.SongDetails, id int) error
	GetText(id, page, size int) ([]string, error)
	GetAllDetails(group, song string) (entity.SongDetails, error)
}

type SongService struct {
	songRepo SongRepository
}

func NewSongService(songRepo SongRepository) *SongService {
	return &SongService{songRepo: songRepo}
}

func (s *SongService) GetSongsPaginated(filter string, page, pageSize int) ([]entity.Song, error) {
	songs, err := s.songRepo.GetAll(filter)
	if err != nil {
		return nil, err
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start > len(songs) {
		start = len(songs)
	}
	if end > len(songs) {
		end = len(songs)
	}

	return songs[start:end], nil

}

func (s *SongService) AddSong(song entity.Song) error {
	var err error
	err = s.songRepo.Add(song)
	if err != nil {
		return err
	}
	return nil
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

func (s *SongService) GetSongLyricsPaginated(id, page, size int) ([]string, error) {
	lyrics, err := s.songRepo.GetText(id, page, size)
	if err != nil {
		return nil, err
	}

	start := page * size
	end := start + size

	if start >= len(lyrics) {
		return []string{}, nil
	}
	if end > len(lyrics) {
		end = len(lyrics)
	}

	return lyrics[start:end], nil
}

func (s *SongService) GetSongInfo(group, song string) (entity.SongDetails, error) {
	var songDetails entity.SongDetails
	var err error
	songDetails, err = s.songRepo.GetAllDetails(group, song)
	if err != nil {
		return songDetails, err
	}
	return songDetails, nil
}
