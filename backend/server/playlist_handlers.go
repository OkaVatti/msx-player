package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetPlaylists(c echo.Context, db *gorm.DB) error {
	var playlists []Playlist
	if err := db.Preload("Songs").Find(&playlists).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db error"})
	}
	return c.JSON(http.StatusOK, playlists)
}

func CreatePlaylistHandler(c echo.Context, db *gorm.DB) error {
	var p Playlist
	if err := c.Bind(&p); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if p.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "name required"})
	}
	if err := db.Create(&p).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db create failed"})
	}
	return c.JSON(http.StatusOK, p)
}

func UpdatePlaylistHandler(c echo.Context, db *gorm.DB) error {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var p Playlist
	if err := db.First(&p, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	var payload Playlist
	if err := c.Bind(&payload); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	p.Name = payload.Name
	p.Description = payload.Description
	// optional: update songs via payload.Songs
	if err := db.Save(&p).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db save failed"})
	}
	return c.JSON(http.StatusOK, p)
}

func DeletePlaylistHandler(c echo.Context, db *gorm.DB) error {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if err := db.Delete(&Playlist{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "delete failed"})
	}
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
