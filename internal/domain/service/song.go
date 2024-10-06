package service

import (
	"EffectiveMobile_Go/internal/domain/entity"
	"go.uber.org/zap"
)

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
	logger   *zap.Logger
}

func NewSongService(songRepo SongRepository, logger *zap.Logger) *SongService {
	return &SongService{songRepo: songRepo, logger: logger}
}

func (s *SongService) GetSongsPaginated(filter string, page, pageSize int) ([]entity.Song, error) {
	s.logger.Debug("Fetching songs paginated", zap.String("filter", filter), zap.Int("page", page), zap.Int("pageSize", pageSize))

	songs, err := s.songRepo.GetAll(filter)
	if err != nil {
		s.logger.Error("Failed to get songs", zap.Error(err))
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

	s.logger.Info("Songs fetched successfully", zap.Int("count", len(songs[start:end])))
	return songs[start:end], nil
}

func (s *SongService) AddSong(song entity.Song) error {
	s.logger.Debug("Adding song", zap.Any("song", song))

	err := s.songRepo.Add(song)
	if err != nil {
		s.logger.Error("Failed to add song", zap.Error(err))
		return err
	}
	s.logger.Info("Song added successfully", zap.String("song", song.Song), zap.String("group", song.Group))
	return nil
}

func (s *SongService) DeleteSong(group string, song string, id int) error {
	s.logger.Debug("Deleting song", zap.String("group", group), zap.String("song", song), zap.Int("id", id))

	err := s.songRepo.Delete(group, song, id)
	if err != nil {
		s.logger.Error("Failed to delete song", zap.Error(err))
		return err
	}
	s.logger.Info("Song deleted successfully", zap.String("group", group), zap.String("song", song))
	return nil
}

func (s *SongService) UpdateSong(song entity.SongDetails, id int) error {
	s.logger.Debug("Updating song", zap.Any("song", song), zap.Int("id", id))

	err := s.songRepo.Update(song, id)
	if err != nil {
		s.logger.Error("Failed to update song", zap.Error(err))
		return err
	}
	s.logger.Info("Song updated successfully", zap.Int("id", id))
	return nil
}

func (s *SongService) GetSongLyricsPaginated(id, page, size int) ([]string, error) {
	s.logger.Debug("Fetching song lyrics paginated", zap.Int("id", id), zap.Int("page", page), zap.Int("size", size))

	lyrics, err := s.songRepo.GetText(id, page, size)
	if err != nil {
		s.logger.Error("Failed to get song lyrics", zap.Error(err))
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

	s.logger.Info("Lyrics fetched successfully", zap.Int("count", len(lyrics[start:end])))
	return lyrics[start:end], nil
}

func (s *SongService) GetSongInfo(group, song string) (entity.SongDetails, error) {
	s.logger.Debug("Fetching song info", zap.String("group", group), zap.String("song", song))

	var songDetails entity.SongDetails
	var err error
	songDetails, err = s.songRepo.GetAllDetails(group, song)
	if err != nil {
		s.logger.Error("Failed to get song info", zap.Error(err))
		return songDetails, err
	}
	s.logger.Info("Song info fetched successfully", zap.Any("songDetails", songDetails))
	return songDetails, nil
}
