package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB
	upgrader    = websocket.Upgrader{}
	clients     = make(map[*websocket.Conn]bool)
	broadcast   = make(chan PlayerState)
	currentSong *Song
	playerState = PlayerState{
		IsPlaying:   false,
		Volume:      0.7,
		Speed:       1.0,
		CurrentTime: 0,
	}
)

type Song struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	Title      string    `json:"title"`
	Artist     string    `json:"artist"`
	Album      string    `json:"album"`
	Genre      string    `json:"genre"`
	Year       int       `json:"year"`
	Duration   float64   `json:"duration"`
	FilePath   string    `json:"file_path" gorm:"unique"`
	Rating     int       `json:"rating"`
	PlayCount  int       `json:"play_count"`
	LastPlayed time.Time `json:"last_played"`
	CreatedAt  time.Time `json:"created_at"`
	Explicit   bool      `json:"explicit"`
	Clean      bool      `json:"clean"`
}

type Playlist struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Songs       []Song    `json:"songs" gorm:"many2many:playlist_songs;"`
	CreatedAt   time.Time `json:"created_at"`
}

type PlayerState struct {
	IsPlaying   bool    `json:"is_playing"`
	CurrentSong *Song   `json:"current_song"`
	CurrentTime float64 `json:"current_time"`
	Duration    float64 `json:"duration"`
	Volume      float64 `json:"volume"`
	Speed       float64 `json:"speed"`
}

func main() {
	// Initialize database
	var err error
	db, err = gorm.Open(sqlite.Open("data/music.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	db.AutoMigrate(&Song{}, &Playlist{})

	// Create music directory if it doesn't exist
	os.MkdirAll("music", 0755)
	os.MkdirAll("data", 0755)

	e := echo.New()

	// Middleware
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/ws", handleWebSocket)
	e.GET("/api/songs", getSongs)
	e.POST("/api/songs", addSong)
	e.DELETE("/api/songs/:id", deleteSong)
	e.PUT("/api/songs/:id", updateSong)
	e.GET("/api/playlists", getPlaylists)
	e.POST("/api/playlists", createPlaylist)
	e.PUT("/api/playlists/:id", updatePlaylist)
	e.DELETE("/api/playlists/:id", deletePlaylist)
	e.POST("/api/player/play", playSong)
	e.POST("/api/player/pause", pausePlayer)
	e.POST("/api/player/stop", stopPlayer)
	e.POST("/api/player/next", nextSong)
	e.POST("/api/player/previous", previousSong)
	e.POST("/api/player/seek", seekSong)
	e.POST("/api/player/volume", setVolume)
	e.POST("/api/player/speed", setSpeed)
	e.POST("/api/player/rate", rateSong)
	e.GET("/api/player/state", getPlayerState)
	e.Static("/music", "music")

	// Start broadcaster
	go broadcastPlayerState()

	// Scan for music files on startup
	go scanMusicLibrary()

	e.Logger.Fatal(e.Start(":1323"))
}

func handleWebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	clients[ws] = true

	// Send current state to new client
	if err := ws.WriteJSON(playerState); err != nil {
		delete(clients, ws)
		return err
	}

	for {
		var msg map[string]interface{}
		if err := ws.ReadJSON(&msg); err != nil {
			delete(clients, ws)
			break
		}
	}
	return nil
}

func broadcastPlayerState() {
	for {
		state := <-broadcast
		for client := range clients {
			err := client.WriteJSON(state)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func getSongs(c echo.Context) error {
	var songs []Song
	query := db.Model(&Song{})

	// Search
	if search := c.QueryParam("search"); search != "" {
		search = "%" + search + "%"
		query = query.Where("title LIKE ? OR artist LIKE ? OR album LIKE ?", search, search, search)
	}

	// Filter by genre
	if genre := c.QueryParam("genre"); genre != "" {
		query = query.Where("genre = ?", genre)
	}

	// Filter by rating
	if rating := c.QueryParam("rating"); rating != "" {
		query = query.Where("rating >= ?", rating)
	}

	// Sort
	sortBy := c.QueryParam("sort_by")
	if sortBy == "" {
		sortBy = "title"
	}
	query = query.Order(sortBy)

	if err := query.Find(&songs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, songs)
}

func addSong(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
	}

	// Check if file is MP3
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".mp3") {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Only MP3 files are supported"})
	}

	// Save file
	filePath := filepath.Join("music", file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Extract metadata and create song record
	song := Song{
		Title:     strings.TrimSuffix(file.Filename, ".mp3"),
		Artist:    "Unknown Artist",
		Album:     "Unknown Album",
		FilePath:  filePath,
		Rating:    0,
		PlayCount: 0,
		CreatedAt: time.Now(),
	}

	// Here you would extract actual metadata from MP3 file
	// For now, we'll use filename as title

	// Check for duplicates
	var existingSong Song
	if err := db.Where("file_path = ?", filePath).First(&existingSong).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Song already exists"})
	}

	// Remove duplicates based on title and artist
	var duplicateSongs []Song
	db.Where("title = ? AND artist = ?", song.Title, song.Artist).Find(&duplicateSongs)

	if len(duplicateSongs) > 0 {
		// Keep the highest rated or oldest
		sort.Slice(duplicateSongs, func(i, j int) bool {
			if duplicateSongs[i].Rating != duplicateSongs[j].Rating {
				return duplicateSongs[i].Rating > duplicateSongs[j].Rating
			}
			return duplicateSongs[i].ID < duplicateSongs[j].ID
		})

		// Delete all except the first one
		for i := 1; i < len(duplicateSongs); i++ {
			db.Delete(&duplicateSongs[i])
		}
	}

	if err := db.Create(&song).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, song)
}

// Implement other handler functions...
// [Previous imports and types remain the same...]

func deleteSong(c echo.Context) error {
	id := c.Param("id")
	if err := db.Delete(&Song{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Song deleted"})
}

func updateSong(c echo.Context) error {
	id := c.Param("id")
	var song Song
	if err := db.First(&song, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Song not found"})
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	if err := db.Model(&song).Updates(updates).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, song)
}

func playSong(c echo.Context) error {
	var request struct {
		SongID uint `json:"song_id"`
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	var song Song
	if err := db.First(&song, request.SongID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Song not found"})
	}

	currentSong = &song
	playerState.IsPlaying = true
	playerState.CurrentSong = currentSong
	playerState.CurrentTime = 0

	// Update play count
	db.Model(&song).Update("play_count", song.PlayCount+1)
	db.Model(&song).Update("last_played", time.Now())

	broadcast <- playerState
	return c.JSON(http.StatusOK, playerState)
}

func pausePlayer(c echo.Context) error {
	playerState.IsPlaying = false
	broadcast <- playerState
	return c.JSON(http.StatusOK, playerState)
}

func stopPlayer(c echo.Context) error {
	playerState.IsPlaying = false
	playerState.CurrentTime = 0
	broadcast <- playerState
	return c.JSON(http.StatusOK, playerState)
}

func seekSong(c echo.Context) error {
	var request struct {
		Time float64 `json:"time"`
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	playerState.CurrentTime = request.Time
	broadcast <- playerState
	return c.JSON(http.StatusOK, playerState)
}

func setVolume(c echo.Context) error {
	var request struct {
		Volume float64 `json:"volume"`
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	playerState.Volume = request.Volume
	broadcast <- playerState
	return c.JSON(http.StatusOK, playerState)
}

func setSpeed(c echo.Context) error {
	var request struct {
		Speed float64 `json:"speed"`
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	playerState.Speed = request.Speed
	broadcast <- playerState
	return c.JSON(http.StatusOK, playerState)
}

func rateSong(c echo.Context) error {
	var request struct {
		SongID uint `json:"song_id"`
		Rating int  `json:"rating"`
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	var song Song
	if err := db.First(&song, request.SongID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Song not found"})
	}

	if request.Rating < 0 || request.Rating > 10 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Rating must be between 0 and 10"})
	}

	db.Model(&song).Update("rating", request.Rating)

	return c.JSON(http.StatusOK, map[string]string{"message": "Rating updated"})
}

func getPlayerState(c echo.Context) error {
	return c.JSON(http.StatusOK, playerState)
}

func scanMusicLibrary() {
	files, err := filepath.Glob("music/*.mp3")
	if err != nil {
		log.Println("Error scanning music library:", err)
		return
	}

	for _, file := range files {
		var existingSong Song
		if err := db.Where("file_path = ?", file).First(&existingSong).Error; err != nil {
			song := Song{
				Title:     strings.TrimSuffix(filepath.Base(file), ".mp3"),
				Artist:    "Unknown Artist",
				Album:     "Unknown Album",
				FilePath:  file,
				Rating:    0,
				PlayCount: 0,
				CreatedAt: time.Now(),
			}
			db.Create(&song)
		}
	}
}

// Playlist handlers
func getPlaylists(c echo.Context) error {
	var playlists []Playlist
	db.Preload("Songs").Find(&playlists)
	return c.JSON(http.StatusOK, playlists)
}

func createPlaylist(c echo.Context) error {
	var playlist Playlist
	if err := json.NewDecoder(c.Request().Body).Decode(&playlist); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	playlist.CreatedAt = time.Now()
	if err := db.Create(&playlist).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, playlist)
}
