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
	"strconv"
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

	// Получаем фильтр
	filter := r.URL.Query().Get("filter")

	// Получаем все песни с учетом фильтра
	songs, err = a.serv.GetSongs(filter)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Получаем параметры пагинации
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	// Устанавливаем значения по умолчанию
	page := 1
	pageSize := 5

	// Преобразуем параметры из строки в int
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

	// Расчет индексов для среза
	start := (page - 1) * pageSize
	end := start + pageSize

	// Обеспечиваем корректные границы
	if start > len(songs) {
		start = len(songs)
	}
	if end > len(songs) {
		end = len(songs)
	}

	// Срез песен в соответствии с пагинацией
	paginatedSongs := songs[start:end]

	// Установка заголовка
	w.Header().Set("Content-Type", "application/json")

	// Кодируем результат в JSON и отправляем
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
	// http.Error(w, "Ok", http.StatusOK)
	resp, err := http.Get("http://localhost:8080/song/info?group=" + song.Group + "&song=" + song.Song)
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
	/*w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song.ID)*/

}

func (a *SongApp) InfoSongHandler(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")

	if group == "" || song == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

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
	r.HandleFunc("/songs/info", a.InfoSongHandler).Methods("POST")
	// r.HandleFunc("/songs/delete", a.GetSongsHandler).Methods("GET")

	log.Println("Starting HTTP server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
