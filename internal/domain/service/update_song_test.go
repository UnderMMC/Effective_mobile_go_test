package service_test

import (
	"EffectiveMobile_Go/internal/domain/entity"
	mockprovider "EffectiveMobile_Go/internal/domain/repository/mocks"
	"EffectiveMobile_Go/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type testDataGenerateUpdateSong = func(
	mp *mockprovider.SongRepository,
	t *testing.T,
) *testDataUpdateSong

type testDataUpdateSong struct {
	name        string
	req         entity.SongDetails
	expectedErr error
}

func successeDataUpdateSong(
	mp *mockprovider.SongRepository,
	t *testing.T,
) *testDataUpdateSong {
	req := entity.SongDetails{
		ReleaseDate: "8 March",
		Text:        "Hello",
		Link:        "http:/example.com",
	}

	//При передаче первоначальной структуры так делать не надо, но мы так делаем при проверке смапленной структуры
	mp.On("Update", mock.MatchedBy(func(val *entity.SongDetails) bool {
		return assert.Equal(t, val.ReleaseDate, "8 March") &&
			assert.Equal(t, val.Text, "Hello") &&
			assert.Equal(t, val.Link, "http:/example.com")
	}), 1).Return(nil)

	return &testDataUpdateSong{
		name:        "successeDataUpdateSong",
		req:         req,
		expectedErr: nil,
	}
}

func providerErrorDataUpdateSong(
	mp *mockprovider.SongRepository,
	t *testing.T,
) *testDataUpdateSong {
	req := entity.SongDetails{
		ReleaseDate: "8 March",
		Text:        "Hello",
		Link:        "http:/example.com",
	}

	mp.On("Update", req, 1).Return(assert.AnError)

	return &testDataUpdateSong{
		name:        "providerErrorDataUpdateSong",
		req:         req,
		expectedErr: assert.AnError,
	}
}

func TestUpdateSong(t *testing.T) {
	testDataFn := []testDataGenerateUpdateSong{
		successeDataUpdateSong,
		providerErrorDataUpdateSong,
	}

	for _, fn := range testDataFn {
		mockProvider := new(mockprovider.SongRepository)

		testArgs := fn(mockProvider, t)

		useCase := service.NewSongService(mockProvider, nil)

		t.Run(testArgs.name, func(t *testing.T) {
			err := useCase.UpdateSong(testArgs.req, 1)
			assert.Equal(t, testArgs.expectedErr, err)
		})

		assert.True(t, mockProvider.AssertExpectations(t))
	}
}
