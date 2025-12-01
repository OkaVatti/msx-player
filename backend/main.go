package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db           *gorm.DB
	upgrader     = websocket.Upgrader{}
	clients      = make(map[*websocket.Conn]bool)
	broadcast    = make(chan PlayerState)
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

type PlaylistSong struct {
	PlaylistID uint `gorm:"primaryKey"`
	SongID     uint `gorm:"primaryKey"`
	Order      int  `json:"order"`
}

type PlayerState struct {
	IsPlaying   bool    `json:"is_playing"`
	CurrentSong *Song   `json:"current_song"`
	CurrentTime float64 `json:"current_time"`
	Duration    float64 `json:"duration"`
	Volume      float64 `json:"volume"`
	Speed       float64 `json:"speed"`
}

type UploadResponse struct {
	Message string `json:"message"`
	Song    *Song  `json:"song,omitempty"`
	Error   string `json:"error,omitempty"`
}

func init() {
	// Create directories on startup
	if err := os.MkdirAll("music", 0755); err != nil {
		log.Fatalf("Failed to create music directory: %v", err)
	}
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}
}

func main() {
	var err error

	// Initialize database with absolute path
	dbPath := filepath.Join("data", "music.db")
	fmt.Printf("📁 Database path: %s\n", dbPath)

	// Check if data directory exists
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		log.Fatal("Data directory does not exist. Please create it manually or check permissions.")
	}

	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("✅ Database connection established")

	// Auto migrate
	if err := db.AutoMigrate(&Song{}, &Playlist{}, &PlaylistSong{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("✅ Database migration completed")

	e := echo.New()

	// Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:1323"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("50M")) // Increase for large MP3 files

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "MSX Player API Server",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	e.GET("/ws", handleWebSocket)
	e.GET("/api/songs", getSongs)
	e.POST("/api/songs", addSong)
	e.DELETE("/api/songs/:id", deleteSong)
	e.PUT("/api/songs/:id", updateSong)
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

	// Static file serving with proper headers
	e.Static("/music", "music")
	e.GET("/music/*", func(c echo.Context) error {
		filename := c.Param("*")
		safePath := filepath.Join("music", filepath.Clean("/"+filename))

		// Check if file exists
		if _, err := os.Stat(safePath); os.IsNotExist(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "File not found"})
		}

		// Set proper headers for audio files
		c.Response().Header().Set("Content-Type", "audio/mpeg")
		c.Response().Header().Set("Cache-Control", "public, max-age=3600")

		return c.File(safePath)
	})

	// Start broadcaster
	go broadcastPlayerState()

	// Scan for music files on startup
	go scanMusicLibrary()

	fmt.Println("🚀 MSX Player Backend started on :1323")
	fmt.Println("📁 Music directory: ./music")
	fmt.Println("💾 Database: ./data/music.db")
	fmt.Println("🌐 API available at: http://localhost:1323")
	fmt.Println("🔌 WebSocket available at: ws://localhost:1323/ws")

	e.Logger.Fatal(e.Start(":1323"))
}

// [Rest of the functions remain the same as in the previous version]
// handleWebSocket, broadcastPlayerState, getSongs, addSong, etc.
// ... (include all the other functions from the previous version)

func handleWebSocket(c echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
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
		// Handle incoming WebSocket messages if needed
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
		search = "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(artist) LIKE ? OR LOWER(album) LIKE ?", search, search, search)
	}

	// Filter by genre
	if genre := c.QueryParam("genre"); genre != "" {
		query = query.Where("genre = ?", genre)
	}

	// Filter by rating
	if rating := c.QueryParam("rating"); rating != "" {
		if ratingInt, err := strconv.Atoi(rating); err == nil {
			query = query.Where("rating >= ?", ratingInt)
		}
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
		return c.JSON(http.StatusBadRequest, UploadResponse{Error: "No file uploaded"})
	}

	// Check if file is MP3
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".mp3") {
		return c.JSON(http.StatusBadRequest, UploadResponse{Error: "Only MP3 files are supported"})
	}

	// Generate safe filename
	safeFilename := generateSafeFilename(file.Filename)
	filePath := filepath.Join("music", safeFilename)

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		return c.JSON(http.StatusConflict, UploadResponse{Error: "File already exists"})
	}

	// Save file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, UploadResponse{Error: "Failed to open uploaded file"})
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, UploadResponse{Error: "Failed to create file on server"})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, UploadResponse{Error: "Failed to save file"})
	}

	// Extract basic metadata from filename
	title, artist, album := extractMetadataFromFilename(file.Filename)

	// Create song record
	song := Song{
		Title:      title,
		Artist:     artist,
		Album:      album,
		Genre:      "Unknown",
		Year:       time.Now().Year(),
		Duration:   0, // Would need MP3 parsing to get actual duration
		FilePath:   filePath,
		Rating:     0,
		PlayCount:  0,
		LastPlayed: time.Time{},
		CreatedAt:  time.Now(),
		Explicit:   strings.Contains(strings.ToLower(file.Filename), "explicit"),
		Clean:      strings.Contains(strings.ToLower(file.Filename), "clean"),
	}

	// Check for duplicates and handle according to rules
	if err := handleDuplicateSongs(&song); err != nil {
		// Remove the uploaded file if duplicate handling fails
		os.Remove(filePath)
		return c.JSON(http.StatusConflict, UploadResponse{Error: err.Error()})
	}

	if err := db.Create(&song).Error; err != nil {
		// Remove the uploaded file if database creation fails
		os.Remove(filePath)
		return c.JSON(http.StatusInternalServerError, UploadResponse{Error: "Failed to create song record"})
	}

	return c.JSON(http.StatusCreated, UploadResponse{
		Message: "Song uploaded successfully",
		Song:    &song,
	})
}

func handleDuplicateSongs(newSong *Song) error {
	// Check for exact file path duplicates
	var existingByPath Song
	if err := db.Where("file_path = ?", newSong.FilePath).First(&existingByPath).Error; err == nil {
		return fmt.Errorf("song with this file path already exists")
	}

	// Check for content duplicates (same title and artist)
	var duplicates []Song
	if err := db.Where("title = ? AND artist = ?", newSong.Title, newSong.Artist).Find(&duplicates).Error; err != nil {
		return nil // No duplicates found
	}

	if len(duplicates) == 0 {
		return nil // No duplicates
	}

	// Check if we have explicit/clean version exception
	hasExplicitCleanException := false
	for _, dup := range duplicates {
		if (newSong.Explicit && dup.Clean) || (newSong.Clean && dup.Explicit) {
			hasExplicitCleanException = true
			break
		}
	}

	if !hasExplicitCleanException {
		// Remove duplicates according to rules: highest rating, then oldest ID
		sort.Slice(duplicates, func(i, j int) bool {
			if duplicates[i].Rating != duplicates[j].Rating {
				return duplicates[i].Rating > duplicates[j].Rating
			}
			return duplicates[i].ID < duplicates[j].ID
		})

		// Keep only the best one, delete others
		for i := 1; i < len(duplicates); i++ {
			if err := db.Delete(&duplicates[i]).Error; err != nil {
				return fmt.Errorf("failed to remove duplicate song")
			}
		}
	}

	return nil
}

func generateSafeFilename(filename string) string {
	// Remove path components and special characters
	base := filepath.Base(filename)
	base = strings.ReplaceAll(base, "..", "")
	base = strings.ReplaceAll(base, "/", "")
	base = strings.ReplaceAll(base, "\\", "")

	// Add timestamp to make unique
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	timestamp := time.Now().UnixNano()

	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}

func extractMetadataFromFilename(filename string) (title, artist, album string) {
	base := strings.TrimSuffix(filepath.Base(filename), ".mp3")

	// Simple parsing: assume "Artist - Title" or "Artist - Title - Album" format
	parts := strings.Split(base, " - ")

	switch len(parts) {
	case 1:
		title = parts[0]
		artist = "Unknown Artist"
		album = "Unknown Album"
	case 2:
		artist = parts[0]
		title = parts[1]
		album = "Unknown Album"
	default:
		artist = parts[0]
		title = parts[1]
		album = strings.Join(parts[2:], " - ")
	}

	// Clean up any explicit/clean tags from titles
	title = strings.ReplaceAll(title, "(Explicit)", "")
	title = strings.ReplaceAll(title, "(Clean)", "")
	title = strings.ReplaceAll(title, "[Explicit]", "")
	title = strings.ReplaceAll(title, "[Clean]", "")
	title = strings.TrimSpace(title)

	return title, artist, album
}

func deleteSong(c echo.Context) error {
	id := c.Param("id")

	var song Song
	if err := db.First(&song, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Song not found"})
	}

	// Remove from all playlists first
	db.Exec("DELETE FROM playlist_songs WHERE song_id = ?", id)

	if err := db.Delete(&song).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Note: We don't delete the actual file, just the database record
	return c.JSON(http.StatusOK, map[string]string{"message": "Song removed from library"})
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

// Update the playSong function to handle queue
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

	// If this is a new song (not the current one), reset position
	if currentSong == nil || currentSong.ID != song.ID {
		playerState.CurrentTime = 0
	}

	currentSong = &song
	playerState.IsPlaying = true
	playerState.CurrentSong = currentSong
	playerState.Duration = song.Duration

	// Update play count and last played
	db.Model(&song).Updates(map[string]interface{}{
		"play_count":  song.PlayCount + 1,
		"last_played": time.Now(),
	})

	broadcast <- playerState
	return c.JSON(http.StatusOK, playerState)
}

func shuffleSongs(c echo.Context) error {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	var songs []Song
	if err := db.Order("title").Find(&songs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get songs"})
	}

	// Convert to pointers and shuffle
	currentQueue = make([]*Song, len(songs))
	for i := range songs {
		currentQueue[i] = &songs[i]
	}

	// Fisher-Yates shuffle
	rand.Shuffle(len(currentQueue), func(i, j int) {
		currentQueue[i], currentQueue[j] = currentQueue[j], currentQueue[i]
	})

	currentIndex = -1
	return c.JSON(http.StatusOK, map[string]string{"message": "Queue shuffled"})
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

func nextSong(c echo.Context) error {
	queueMutex.RLock()
	defer queueMutex.RUnlock()

	if len(currentQueue) == 0 {
		// If no queue, get all songs
		var songs []Song
		if err := db.Order("title").Find(&songs).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get songs"})
		}

		// Convert to pointers
		currentQueue = make([]*Song, len(songs))
		for i := range songs {
			currentQueue[i] = &songs[i]
		}
	}

	if len(currentQueue) > 0 {
		if currentIndex < len(currentQueue)-1 {
			currentIndex++
		} else {
			currentIndex = 0 // Loop back to start
		}

		currentSong = currentQueue[currentIndex]
		playerState.IsPlaying = true
		playerState.CurrentSong = currentSong
		playerState.CurrentTime = 0
		playerState.Duration = currentSong.Duration

		// Update play count
		db.Model(currentSong).Updates(map[string]interface{}{
			"play_count":  currentSong.PlayCount + 1,
			"last_played": time.Now(),
		})

		broadcast <- playerState
	}

	return c.JSON(http.StatusOK, playerState)
}

// Update previousSong to actually play previous song
func previousSong(c echo.Context) error {
	queueMutex.RLock()
	defer queueMutex.RUnlock()

	if len(currentQueue) > 0 {
		if currentIndex > 0 {
			currentIndex--
		} else {
			currentIndex = len(currentQueue) - 1 // Loop to end
		}

		currentSong = currentQueue[currentIndex]
		playerState.IsPlaying = true
		playerState.CurrentSong = currentSong
		playerState.CurrentTime = 0
		playerState.Duration = currentSong.Duration

		// Update play count
		db.Model(currentSong).Updates(map[string]interface{}{
			"play_count":  currentSong.PlayCount + 1,
			"last_played": time.Now(),
		})

		broadcast <- playerState
	}

	return c.JSON(http.StatusOK, playerState)
}

func scanMusicLibrary() {
	fmt.Println("🔍 Scanning music library...")

	files, err := filepath.Glob("music/*.mp3")
	if err != nil {
		log.Println("Error scanning music library:", err)
		return
	}

	for _, file := range files {
		var existingSong Song
		if err := db.Where("file_path = ?", file).First(&existingSong).Error; err != nil {
			// File not in database, add it
			title, artist, album := extractMetadataFromFilename(filepath.Base(file))

			song := Song{
				Title:     title,
				Artist:    artist,
				Album:     album,
				FilePath:  file,
				Rating:    0,
				PlayCount: 0,
				CreatedAt: time.Now(),
				Explicit:  strings.Contains(strings.ToLower(file), "explicit"),
				Clean:     strings.Contains(strings.ToLower(file), "clean"),
			}

			// Handle duplicates for scanned files too
			if err := handleDuplicateSongs(&song); err == nil {
				db.Create(&song)
				fmt.Printf("➕ Added: %s - %s\n", song.Artist, song.Title)
			}
		}
	}
	fmt.Printf("✅ Music library scan complete. Found %d files.\n", len(files))
}

// Playlist handlers
func getPlaylists(c echo.Context) error {
	var playlists []Playlist
	db.Preload("Songs").Find(&playlists)
	return c.JSON(http.StatusOK, playlists)
}

func createPlaylist(c echo.Context) error {
	var playlist struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&playlist); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	newPlaylist := Playlist{
		Name:        playlist.Name,
		Description: playlist.Description,
		CreatedAt:   time.Now(),
	}

	if err := db.Create(&newPlaylist).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, newPlaylist)
}

func updatePlaylist(c echo.Context) error {
	id := c.Param("id")
	var playlist Playlist
	if err := db.First(&playlist, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Playlist not found"})
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(c.Request().Body).Decode(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	if err := db.Model(&playlist).Updates(updates).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, playlist)
}

func deletePlaylist(c echo.Context) error {
	id := c.Param("id")

	var playlist Playlist
	if err := db.First(&playlist, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Playlist not found"})
	}

	// Clear playlist songs association first
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

	// Add song to playlist
	if err := db.Model(&playlist).Association("Songs").Append(&song); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Song added to playlist"})
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

	// Remove song from playlist
	if err := db.Model(&playlist).Association("Songs").Delete(&song); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Song removed from playlist"})
}
