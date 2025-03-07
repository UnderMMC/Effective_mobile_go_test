package service

import (
	"EffectiveMobile_Go/internal/domain/entity"
	"go.uber.org/zap"
)

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
