package service

import "EffectiveMobile_Go/internal/domain/entity"

type SongRepository interface {
	GetAll(filter string) ([]entity.Song, error)
	//GetByID(id int) (entity.Song, error)
	//Add(song entity.Song) (int, error)
	//Update(song entity.Song) error
	//Delete(id int) error
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

/*func (s *SongService) AddSong(song entity.Song) (int, error) {
	// Здесь можно добавить валидацию данных песни перед добавлением
	id, err := s.songRepo.Add(song)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *SongService) UpdateSong(song entity.Song) error {
	// Здесь можно добавить валидацию данных песни перед обновлением
	err := s.songRepo.Update(song)
	if err != nil {
		return err
	}
	return nil
}

func (s *SongService) DeleteSong(id int) error {
	err := s.songRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}*/
