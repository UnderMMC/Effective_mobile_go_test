package app

import (
	"EffectiveMobile_Go/internal/domain/entity"
	"EffectiveMobile_Go/internal/domain/repository"
	"EffectiveMobile_Go/internal/domain/service"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var db *sql.DB

type SongService interface {
	GetSongs(filter string) ([]entity.Song, error)
	AddSong(song entity.Song) (int, error)
	DeleteSong(group string, song string, id int) error
	UpdateSong(song entity.SongDetails, id int) error
	GetSongLyricsPaginated(id, page, size int) ([]string, error)
}

type SongApp struct {
	serv SongService
}

var songs = map[string]entity.SongDetails{
	"Supermassive Black Hole": {
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	},
}

func (a *SongApp) GetSongsHandler(w http.ResponseWriter, r *http.Request) {
	var songs []entity.Song
	var err error

	filter := r.URL.Query().Get("filter")

	songs, err = a.serv.GetSongs(filter)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page := 1
	pageSize := 5

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start > len(songs) {
		start = len(songs)
	}
	if end > len(songs) {
		end = len(songs)
	}

	paginatedSongs := songs[start:end]

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(paginatedSongs); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
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

	if song.Group == "" || song.Song == "" {
		http.Error(w, "Group or Song cannot be empty", http.StatusBadRequest)
		return
	}

	// Кодирование параметров для URL
	groupEncoded := url.QueryEscape(song.Group)
	songEncoded := url.QueryEscape(song.Song)

	resp, err := http.Get("http://localhost:8080/songs/info?group=" + groupEncoded + "&song=" + songEncoded)
	if err != nil {
		http.Error(w, "Failed to fetch song info", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (a *SongApp) InfoSongHandler(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")

	if group == "" || song == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли песня в карте.
	songDetail, exists := songs[song]
	if !exists {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(songDetail); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (a *SongApp) DeleteSongHandler(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	if group == "" || song == "" || id == 0 || id < 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	err := a.serv.DeleteSong(group, song, id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
	}
}

func (a *SongApp) UpdateSongHandler(w http.ResponseWriter, r *http.Request) {
	var song entity.SongDetails
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	if id == 0 || id < 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = a.serv.UpdateSong(song, id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Song updated successfully"))
}

func (a *SongApp) GetTextHandler(w http.ResponseWriter, r *http.Request) {
	//var err error
	idStr := r.URL.Query().Get("id")
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")

	id, _ := strconv.Atoi(idStr)
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	if id <= 0 || page < 0 || size <= 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	//var text []string
	text, err := a.serv.GetSongLyricsPaginated(id, page, size)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(text); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
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
	r.HandleFunc("/songs/info", a.InfoSongHandler).Methods("GET")
	r.HandleFunc("/songs/delete", a.DeleteSongHandler).Methods("GET")
	r.HandleFunc("/songs/update", a.UpdateSongHandler).Methods("POST")
	r.HandleFunc("/songs/text", a.GetTextHandler).Methods("GET")

	log.Println("Starting HTTP server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
