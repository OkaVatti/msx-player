package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// simple in-memory player state for the server
var globalPlayerState = struct {
	IsPlaying   bool    `json:"is_playing"`
	CurrentSong *Song   `json:"current_song"`
	CurrentTime float64 `json:"current_time"`
	Duration    float64 `json:"duration"`
	Volume      float64 `json:"volume"`
	Speed       float64 `json:"speed"`
}{
	IsPlaying:   false,
	CurrentSong: nil,
	CurrentTime: 0,
	Duration:    0,
	Volume:      0.7,
	Speed:       1.0,
}

func PlayerStateHandler(c echo.Context, db *gorm.DB) error {
	return c.JSON(http.StatusOK, globalPlayerState)
}

type controlPayload struct {
	SongID int     `json:"song_id"`
	Time   float64 `json:"time"`
	Rating int     `json:"rating"`
}

func PlayerPlayHandler(c echo.Context, db *gorm.DB) error {
	var p controlPayload
	if err := c.Bind(&p); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if p.SongID == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "song_id required"})
	}
	var s Song
	if err := db.First(&s, p.SongID).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	globalPlayerState.CurrentSong = &s
	globalPlayerState.Duration = s.Duration
	globalPlayerState.IsPlaying = true
	globalPlayerState.CurrentTime = 0
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func PlayerPauseHandler(c echo.Context, db *gorm.DB) error {
	globalPlayerState.IsPlaying = false
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func PlayerNextHandler(c echo.Context, db *gorm.DB) error {
	// For demo, set playing false — client should manage queue
	globalPlayerState.IsPlaying = false
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func PlayerPreviousHandler(c echo.Context, db *gorm.DB) error {
	globalPlayerState.IsPlaying = false
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func PlayerSeekHandler(c echo.Context, db *gorm.DB) error {
	var p controlPayload
	if err := c.Bind(&p); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	globalPlayerState.CurrentTime = p.Time
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func PlayerRateHandler(c echo.Context, db *gorm.DB) error {
	var p controlPayload
	if err := c.Bind(&p); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if p.SongID == 0 {
		return c.NoContent(http.StatusBadRequest)
	}
	var s Song
	if err := db.First(&s, p.SongID).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	s.Rating = p.Rating
	if err := db.Save(&s).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db save failed"})
	}
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
