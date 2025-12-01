package main

import (
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dhowden/tag"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db       *gorm.DB
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients      = make(map[*websocket.Conn]bool)
	clientsMutex sync.RWMutex
	broadcast    = make(chan PlayerState, 100)
	currentSong  *Song
	currentQueue []*Song
	currentIndex int = -1
	queueMutex   sync.RWMutex
	playerState  = PlayerState{
		IsPlaying:   false,
		Volume:      0.7,
		Speed:       1.0,
		CurrentTime: 0,
	}
	stateMutex sync.RWMutex
)

type Song struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	Title      string    `json:"title" gorm:"index"`
	Artist     string    `json:"artist" gorm:"index"`
	Album      string    `json:"album" gorm:"index"`
	Genre      string    `json:"genre" gorm:"index"`
	Year       int       `json:"year"`
	Duration   float64   `json:"duration"`
	FilePath   string    `json:"file_path" gorm:"unique"`
	FileHash   string    `json:"file_hash" gorm:"unique;index"`
	Rating     int       `json:"rating" gorm:"default:0"`
	PlayCount  int       `json:"play_count" gorm:"default:0"`
	LastPlayed time.Time `json:"last_played"`
	CreatedAt  time.Time `json:"created_at"`
	Explicit   bool      `json:"explicit" gorm:"default:false"`
	Clean      bool      `json:"clean" gorm:"default:false"`
	TrackNum   int       `json:"track_number"`
}

type Playlist struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"index"`
	Description string    `json:"description"`
	Songs       []Song    `json:"songs" gorm:"many2many:playlist_songs;"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PlaylistSong struct {
	PlaylistID uint `gorm:"primaryKey"`
	SongID     uint `gorm:"primaryKey"`
	Position   int  `json:"position" gorm:"default:0"`
}

type PlayerState struct {
	IsPlaying   bool    `json:"is_playing"`
	CurrentSong *Song   `json:"current_song"`
	CurrentTime float64 `json:"current_time"`
	Duration    float64 `json:"duration"`
	Volume      float64 `json:"volume"`
	Speed       float64 `json:"speed"`
}

func init() {
	dirs := []string{"music", "data"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create %s directory: %v", dir, err)
		}
	}
}

func main() {
	var err error
	dbPath := filepath.Join("data", "music.db")

	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Song{}, &Playlist{}, &PlaylistSong{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:1323"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("50M"))

	// Routes
	e.GET("/", healthCheck)
	e.GET("/ws", handleWebSocket)
	e.GET("/api/songs", getSongs)
	e.POST("/api/songs", addSong)
	e.DELETE("/api/songs/:id", deleteSong)
	e.PUT("/api/songs/:id", updateSong)
	e.GET("/api/songs/:id/stream", streamSong)
	e.GET("/api/playlists", getPlaylists)
	e.POST("/api/playlists", createPlaylist)
	e.PUT("/api/playlists/:id", updatePlaylist)
	e.DELETE("/api/playlists/:id", deletePlaylist)
	e.POST("/api/playlists/:id/songs/:songId", addSongToPlaylist)
	e.DELETE("/api/playlists/:id/songs/:songId", removeSongFromPlaylist)
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
	e.POST("/api/player/shuffle", shuffleQueue)

	go broadcastPlayerState()
	go scanMusicLibrary()

	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║     MSX PLAYER BACKEND v1.0.0          ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Printf("📁 Music: ./music\n")
	fmt.Printf("💾 Database: %s\n", dbPath)
	fmt.Printf("🌐 API: http://localhost:1323\n")
	fmt.Printf("🔌 WebSocket: ws://localhost:1323/ws\n")

	e.Logger.Fatal(e.Start(":1323"))
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "running",
		"version": "1.0.0",
		"songs":   getSongCount(),
	})
}

func getSongCount() int {
	var count int64
	db.Model(&Song{}).Count(&count)
	return int(count)
}

func handleWebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer func() {
		clientsMutex.Lock()
		delete(clients, ws)
		clientsMutex.Unlock()
		ws.Close()
	}()

	clientsMutex.Lock()
	clients[ws] = true
	clientsMutex.Unlock()

	stateMutex.RLock()
	state := playerState
	stateMutex.RUnlock()

	if err := ws.WriteJSON(state); err != nil {
		return err
	}

	for {
		var msg map[string]interface{}
		if err := ws.ReadJSON(&msg); err != nil {
			break
		}
	}
	return nil
}

func broadcastPlayerState() {
	for state := range broadcast {
		clientsMutex.RLock()
		for client := range clients {
			if err := client.WriteJSON(state); err != nil {
				client.Close()
				clientsMutex.Lock()
				delete(clients, client)
				clientsMutex.Unlock()
			}
		}
		clientsMutex.RUnlock()
	}
}

func getSongs(c echo.Context) error {
	var songs []Song
	query := db.Model(&Song{})

	if search := c.QueryParam("search"); search != "" {
		search = "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(artist) LIKE ? OR LOWER(album) LIKE ?",
			search, search, search)
	}

	if genre := c.QueryParam("genre"); genre != "" {
		query = query.Where("genre = ?", genre)
	}

	if rating := c.QueryParam("rating"); rating != "" {
		if r, err := strconv.Atoi(rating); err == nil {
			query = query.Where("rating >= ?", r)
		}
	}

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

func streamSong(c echo.Context) error {
	id := c.Param("id")
	var song Song

	if err := db.First(&song, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Song not found"})
	}

	file, err := os.Open(song.FilePath)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "File not found"})
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to stat file"})
	}

	c.Response().Header().Set("Content-Type", "audio/mpeg")
	c.Response().Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	c.Response().Header().Set("Accept-Ranges", "bytes")

	http.ServeContent(c.Response(), c.Request(), song.FilePath, stat.ModTime(), file)
	return nil
}

func addSong(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
	}

	if !strings.HasSuffix(strings.ToLower(file.Filename), ".mp3") {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Only MP3 files supported"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open file"})
	}
	defer src.Close()

	// Extract metadata
	metadata, err := tag.ReadFrom(src)
	src.Seek(0, 0)

	song := Song{
		Title:    file.Filename,
		Artist:   "Unknown Artist",
		Album:    "Unknown Album",
		Genre:    "Unknown",
		Year:     time.Now().Year(),
		Duration: 0,
		Rating:   0,
	}

	if metadata != nil {
		if title := metadata.Title(); title != "" {
			song.Title = title
		}
		if artist := metadata.Artist(); artist != "" {
			song.Artist = artist
		}
		if album := metadata.Album(); album != "" {
			song.Album = album
		}
		if genre := metadata.Genre(); genre != "" {
			song.Genre = genre
		}
		if year := metadata.Year(); year != 0 {
			song.Year = year
		}
		trackNum, _ := metadata.Track()
		song.TrackNum = trackNum
	}

	safeFilename := generateSafeFilename(file.Filename)
	filePath := filepath.Join("music", safeFilename)
	song.FilePath = filePath

	dst, err := os.Create(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		os.Remove(filePath)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to write file"})
	}

	if err := db.Create(&song).Error; err != nil {
		os.Remove(filePath)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Song uploaded successfully",
		"song":    song,
	})
}

func generateSafeFilename(filename string) string {
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, name)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}

func deleteSong(c echo.Context) error {
	id := c.Param("id")
	var song Song

	if err := db.First(&song, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Song not found"})
	}

	db.Exec("DELETE FROM playlist_songs WHERE song_id = ?", id)

	if err := db.Delete(&song).Error; err != nil {
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
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid data"})
	}

	if err := db.Model(&song).Updates(updates).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	db.First(&song, id)
	return c.JSON(http.StatusOK, song)
}

func playSong(c echo.Context) error {
	var req struct {
		SongID uint `json:"song_id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	var song Song
	if err := db.First(&song, req.SongID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Song not found"})
	}

	stateMutex.Lock()
	currentSong = &song
	playerState.IsPlaying = true
	playerState.CurrentSong = currentSong
	playerState.Duration = song.Duration
	playerState.CurrentTime = 0
	state := playerState
	stateMutex.Unlock()

	db.Model(&song).Updates(map[string]interface{}{
		"play_count":  song.PlayCount + 1,
		"last_played": time.Now(),
	})

	select {
	case broadcast <- state:
	default:
	}

	return c.JSON(http.StatusOK, state)
}

func pausePlayer(c echo.Context) error {
	stateMutex.Lock()
	playerState.IsPlaying = false
	state := playerState
	stateMutex.Unlock()

	select {
	case broadcast <- state:
	default:
	}

	return c.JSON(http.StatusOK, state)
}

func stopPlayer(c echo.Context) error {
	stateMutex.Lock()
	playerState.IsPlaying = false
	playerState.CurrentTime = 0
	state := playerState
	stateMutex.Unlock()

	select {
	case broadcast <- state:
	default:
	}

	return c.JSON(http.StatusOK, state)
}

func nextSong(c echo.Context) error {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	if len(currentQueue) == 0 {
		var songs []Song
		db.Order("title").Find(&songs)
		currentQueue = make([]*Song, len(songs))
		for i := range songs {
			currentQueue[i] = &songs[i]
		}
	}

	if len(currentQueue) > 0 {
		currentIndex = (currentIndex + 1) % len(currentQueue)
		return playSongByIndex(c, currentIndex)
	}

	return c.JSON(http.StatusOK, playerState)
}

func previousSong(c echo.Context) error {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	if len(currentQueue) > 0 {
		currentIndex--
		if currentIndex < 0 {
			currentIndex = len(currentQueue) - 1
		}
		return playSongByIndex(c, currentIndex)
	}

	return c.JSON(http.StatusOK, playerState)
}

func playSongByIndex(c echo.Context, index int) error {
	if index < 0 || index >= len(currentQueue) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid index"})
	}

	song := currentQueue[index]

	stateMutex.Lock()
	currentSong = song
	playerState.IsPlaying = true
	playerState.CurrentSong = currentSong
	playerState.CurrentTime = 0
	playerState.Duration = song.Duration
	state := playerState
	stateMutex.Unlock()

	db.Model(song).Updates(map[string]interface{}{
		"play_count":  song.PlayCount + 1,
		"last_played": time.Now(),
	})

	select {
	case broadcast <- state:
	default:
	}

	return c.JSON(http.StatusOK, state)
}

func seekSong(c echo.Context) error {
	var req struct {
		Time float64 `json:"time"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	stateMutex.Lock()
	playerState.CurrentTime = req.Time
	state := playerState
	stateMutex.Unlock()

	select {
	case broadcast <- state:
	default:
	}

	return c.JSON(http.StatusOK, state)
}

func setVolume(c echo.Context) error {
	var req struct {
		Volume float64 `json:"volume"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	stateMutex.Lock()
	playerState.Volume = req.Volume
	state := playerState
	stateMutex.Unlock()

	select {
	case broadcast <- state:
	default:
	}

	return c.JSON(http.StatusOK, state)
}

func setSpeed(c echo.Context) error {
	var req struct {
		Speed float64 `json:"speed"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	stateMutex.Lock()
	playerState.Speed = req.Speed
	state := playerState
	stateMutex.Unlock()

	select {
	case broadcast <- state:
	default:
	}

	return c.JSON(http.StatusOK, state)
}

func rateSong(c echo.Context) error {
	var req struct {
		SongID uint `json:"song_id"`
		Rating int  `json:"rating"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.Rating < 0 || req.Rating > 10 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Rating must be 0-10"})
	}

	var song Song
	if err := db.First(&song, req.SongID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Song not found"})
	}

	db.Model(&song).Update("rating", req.Rating)
	return c.JSON(http.StatusOK, map[string]string{"message": "Rating updated"})
}

func getPlayerState(c echo.Context) error {
	stateMutex.RLock()
	state := playerState
	stateMutex.RUnlock()
	return c.JSON(http.StatusOK, state)
}

func shuffleQueue(c echo.Context) error {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	rand.Shuffle(len(currentQueue), func(i, j int) {
		currentQueue[i], currentQueue[j] = currentQueue[j], currentQueue[i]
	})

	currentIndex = 0
	return c.JSON(http.StatusOK, map[string]string{"message": "Queue shuffled"})
}

func scanMusicLibrary() {
	time.Sleep(2 * time.Second)

	files, err := filepath.Glob("music/*.mp3")
	if err != nil {
		log.Println("Error scanning:", err)
		return
	}

	for _, filePath := range files {
		var existing Song
		if err := db.Where("file_path = ?", filePath).First(&existing).Error; err == nil {
			continue
		}

		file, err := os.Open(filePath)
		if err != nil {
			continue
		}

		metadata, _ := tag.ReadFrom(file)
		file.Close()

		song := Song{
			Title:    filepath.Base(filePath),
			Artist:   "Unknown Artist",
			Album:    "Unknown Album",
			Genre:    "Unknown",
			Year:     time.Now().Year(),
			FilePath: filePath,
		}

		if metadata != nil {
			if title := metadata.Title(); title != "" {
				song.Title = title
			}
			if artist := metadata.Artist(); artist != "" {
				song.Artist = artist
			}
			if album := metadata.Album(); album != "" {
				song.Album = album
			}
			if genre := metadata.Genre(); genre != "" {
				song.Genre = genre
			}
			if year := metadata.Year(); year != 0 {
				song.Year = year
			}
		}

		db.Create(&song)
		fmt.Printf("➕ %s - %s\n", song.Artist, song.Title)
	}
}

func getPlaylists(c echo.Context) error {
	var playlists []Playlist
	db.Preload("Songs").Find(&playlists)
	return c.JSON(http.StatusOK, playlists)
}

func createPlaylist(c echo.Context) error {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	playlist := Playlist{
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if err := db.Create(&playlist).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, playlist)
}

func updatePlaylist(c echo.Context) error {
	id := c.Param("id")
	var playlist Playlist

	if err := db.First(&playlist, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Playlist not found"})
	}

	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid data"})
	}

	if err := db.Model(&playlist).Updates(updates).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	db.Preload("Songs").First(&playlist, id)
	return c.JSON(http.StatusOK, playlist)
}

func deletePlaylist(c echo.Context) error {
	id := c.Param("id")
	var playlist Playlist

	if err := db.First(&playlist, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Playlist not found"})
	}

	db.Model(&playlist).Association("Songs").Clear()

	if err := db.Delete(&playlist).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Playlist deleted"})
}

func addSongToPlaylist(c echo.Context) error {
	playlistID := c.Param("id")
	songID := c.Param("songId")

	var playlist Playlist
	if err := db.First(&playlist, playlistID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Playlist not found"})
	}

	var song Song
	if err := db.First(&song, songID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Song not found"})
	}

	if err := db.Model(&playlist).Association("Songs").Append(&song); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Song added"})
}

func removeSongFromPlaylist(c echo.Context) error {
	playlistID := c.Param("id")
	songID := c.Param("songId")

	var playlist Playlist
	if err := db.First(&playlist, playlistID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Playlist not found"})
	}

	var song Song
	if err := db.First(&song, songID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Song not found"})
	}

	if err := db.Model(&playlist).Association("Songs").Delete(&song); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Song removed"})
}
