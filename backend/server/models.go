package server
package main

import (
	"time"

	"gorm.io/gorm"
)

// Song model
type Song struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Title      string         `json:"title"`
	Artist     string         `json:"artist"`
	Album      string         `json:"album"`
	Genre      string         `json:"genre"`
	Year       int            `json:"year"`
	Duration   float64        `json:"duration"` // seconds
	FilePath   string         `json:"file_path"` // /files/<filename>
	FileName   string         `json:"file_name"` // internal filename on disk
	OrigName   string         `json:"orig_name"` // original name
	Rating     int            `json:"rating"`    // 0-10
	PlayCount  int            `json:"play_count"`
	LastPlayed *time.Time     `json:"last_played"`
	Explicit   bool           `json:"explicit"`
	Clean      bool           `json:"clean"`
	Size       int64          `json:"size"`
}

// Playlist model
type Playlist struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Songs       []*Song        `gorm:"many2many:playlist_songs;constraint:OnDelete:CASCADE;" json:"songs"`
}
