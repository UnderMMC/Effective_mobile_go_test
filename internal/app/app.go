package app

import (
	"EffectiveMobile_Go/internal/domain/entity"
	"EffectiveMobile_Go/internal/domain/repository"
	"EffectiveMobile_Go/internal/domain/service"
	"database/sql"
	"encoding/json"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var db *sql.DB

type SongService interface {
	GetSongsPaginated(filter string, page, pageSize int) ([]entity.Song, error)
	AddSong(song entity.Song) error
	DeleteSong(group string, song string, id int) error
	UpdateSong(song entity.SongDetails, id int) error
	GetSongLyricsPaginated(id, page, size int) ([]string, error)
	GetSongInfo(group, song string) (entity.SongDetails, error)
}

type SongApp struct {
	serv SongService
}

func (a *SongApp) GetSongsHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	filter := r.URL.Query().Get("filter")

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

	songs, err := a.serv.GetSongsPaginated(filter, page, pageSize)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(songs); err != nil {
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
	err = a.serv.AddSong(song)
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
	var songDetails entity.SongDetails
	var err error
	songDetails, err = a.serv.GetSongInfo(group, song)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(songDetails); err != nil {
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
	if id <= 0 {
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
	var err error
	connStr := "user=postgres password=pgpwd4habr dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := migrateDB(); err != nil {
		log.Fatalf("Migration failed: %v", err)
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

func migrateDB() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	migrationsDir := currentDir + "/migrations"

	// Создаем новое соединение с базой данных
	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// Создаем новый мигратор
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir,
		"postgres", dbDriver)
	if err != nil {
		return err
	}

	// Выполняем миграции
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
