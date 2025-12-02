package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// ensure storage dir
	if err := os.MkdirAll("./storage/audio", 0o755); err != nil {
		log.Fatal(err)
	}

	// open DB
	db, err := gorm.Open(sqlite.Open("msxplayer.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// automigrate
	if err := db.AutoMigrate(&Song{}, &Playlist{}); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// CORS - allow dev origin; tighten in production
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:3333"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{"*"},
	}))

	api := e.Group("/api")

	// file streaming
	e.GET("/files/:filename", StreamFileHandler)

	// song endpoints
	api.GET("/songs", func(c echo.Context) error { return GetSongs(c, db) })
	api.POST("/songs/upload", func(c echo.Context) error { return UploadSongHandler(c, db) })
	api.GET("/songs/:id", func(c echo.Context) error { return GetSongHandler(c, db) })

	// player endpoints (server-side control; clients can fallback to local)
	api.GET("/player/state", func(c echo.Context) error { return PlayerStateHandler(c, db) })
	api.POST("/player/play", func(c echo.Context) error { return PlayerPlayHandler(c, db) })
	api.POST("/player/pause", func(c echo.Context) error { return PlayerPauseHandler(c, db) })
	api.POST("/player/next", func(c echo.Context) error { return PlayerNextHandler(c, db) })
	api.POST("/player/previous", func(c echo.Context) error { return PlayerPreviousHandler(c, db) })
	api.POST("/player/seek", func(c echo.Context) error { return PlayerSeekHandler(c, db) })
	api.POST("/player/rate", func(c echo.Context) error { return PlayerRateHandler(c, db) })

	// playlist endpoints
	api.GET("/playlists", func(c echo.Context) error { return GetPlaylists(c, db) })
	api.POST("/playlists", func(c echo.Context) error { return CreatePlaylistHandler(c, db) })
	api.PUT("/playlists/:id", func(c echo.Context) error { return UpdatePlaylistHandler(c, db) })
	api.DELETE("/playlists/:id", func(c echo.Context) error { return DeletePlaylistHandler(c, db) })

	// static assets for UI if you want to serve frontend from Go
	// e.Static("/", "dist") // if building frontend into dist

	log.Println("server listening on :1323")
	if err := e.Start(":1323"); err != nil {
		log.Fatal(err)
	}
}
