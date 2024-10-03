package app

import (
	"EffectiveMobile_Go/internal/domain/entity"
	"EffectiveMobile_Go/internal/domain/repository"
	"EffectiveMobile_Go/internal/domain/service"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var db *sql.DB

type SongService interface {
	GetSongs(filter string) ([]entity.Song, error)
	AddSong(song entity.Song) (int, error)
	//UpdateSong(song entity.Song) error
	//DeleteSong(id int) error
}

type SongApp struct {
	serv SongService
}

func (a *SongApp) GetSongsHandler(w http.ResponseWriter, r *http.Request) {
	var songs []entity.Song
	var err error
	filter := r.URL.Query().Get("filter")
	songs, err = a.serv.GetSongs(filter)
	if err != nil {
		log.Println(err)
	}
	for _, song := range songs {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(song)
	}
}

func (a *SongApp) AddSongHandler(w http.ResponseWriter, r *http.Request) {
	var song entity.Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	song.ID, err = a.serv.AddSong(song)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.Error(w, "Ok", http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song.ID)

}

func New() *SongApp {
	return &SongApp{}
}

func (a *SongApp) Run() {
	defer db.Close()
	var err error
	connStr := "user=postgres password=pgpwd4habr dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()

	SongRepo := repository.NewPostgresSongRepository(db)
	SongServ := service.NewSongService(SongRepo)
	a.serv = SongServ

	r.HandleFunc("/songs", a.GetSongsHandler).Methods("GET")
	r.HandleFunc("/songs/add", a.AddSongHandler).Methods("POST")

	log.Println("Starting HTTP server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
