package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dhowden/tag"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const storageDir = "./storage/audio"

// Streaming handler supports Range requests
func StreamFileHandler(c echo.Context) error {
	filename := c.Param("filename")
	if filename == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	path := filepath.Join(storageDir, filepath.Clean(filename))
	f, err := os.Open(path)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	size := fi.Size()
	rangeHeader := c.Request().Header.Get("Range")
	if rangeHeader == "" {
		// full content
		return c.Stream(http.StatusOK, "audio/mpeg", f)
	}

	// parse range header
	var start, end int64
	if strings.HasPrefix(rangeHeader, "bytes=") {
		rangeParts := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
		start, _ = strconv.ParseInt(rangeParts[0], 10, 64)
		if rangeParts[1] != "" {
			end, _ = strconv.ParseInt(rangeParts[1], 10, 64)
		} else {
			end = size - 1
		}
		if end >= size {
			end = size - 1
		}
	} else {
		start = 0
		end = size - 1
	}
	if start < 0 || end < start {
		start = 0
		end = size - 1
	}

	_, err = f.Seek(start, io.SeekStart)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	chunkSize := end - start + 1
	c.Response().Header().Set("Content-Type", "audio/mpeg")
	c.Response().Header().Set("Accept-Ranges", "bytes")
	c.Response().Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))
	c.Response().WriteHeader(http.StatusPartialContent)
	_, err = io.CopyN(c.Response(), f, chunkSize)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

// Get songs
func GetSongs(c echo.Context, db *gorm.DB) error {
	var songs []Song
	if err := db.Order("title asc").Find(&songs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db error"})
	}
	return c.JSON(http.StatusOK, songs)
}

// Get single song
func GetSongHandler(c echo.Context, db *gorm.DB) error {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var s Song
	if err := db.First(&s, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.NoContent(http.StatusNotFound)
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db error"})
	}
	return c.JSON(http.StatusOK, s)
}

// Upload handler
func UploadSongHandler(c echo.Context, db *gorm.DB) error {
	// single file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "file required"})
	}

	// save file & extract metadata
	song, err := saveUploadedFile(fileHeader)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// attempt to extract tags
	if err := extractAndFillMetadata(song, fileHeader); err != nil {
		// log but continue
		fmt.Println("metadata extract error:", err)
	}

	// store DB record
	if err := db.Create(song).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db create failed"})
	}

	return c.JSON(http.StatusOK, song)
}

// helpers for saving uploaded file
func saveUploadedFile(fh *multipart.FileHeader) (*Song, error) {
	src, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	id := uuid.New().String()
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if ext == "" {
		ext = ".mp3"
	}
	internalName := id + ext
	dstPath := filepath.Join(storageDir, internalName)
	out, err := os.Create(dstPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	n, err := io.Copy(out, src)
	if err != nil {
		return nil, err
	}

	s := &Song{
		Title:    strings.TrimSuffix(fh.Filename, ext),
		OrigName: fh.Filename,
		FilePath: "/files/" + internalName,
		FileName: internalName,
		Size:     n,
	}
	// attempt to get duration via ffprobe, best-effort
	if dur, derr := probeDuration(dstPath); derr == nil {
		s.Duration = dur
	}
	return s, nil
}

func probeDuration(path string) (float64, error) {
	// ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 file
	out, err := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", path).Output()
	if err != nil {
		return 0, err
	}
	s := strings.TrimSpace(string(out))
	if s == "" {
		return 0, nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

func extractAndFillMetadata(s *Song, fh *multipart.FileHeader) error {
	f, err := fh.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	// tag library reads metadata
	meta, err := tag.ReadFrom(f)
	if err != nil {
		return err
	}
	if meta.Title() != "" {
		s.Title = meta.Title()
	}
	if meta.Artist() != "" {
		s.Artist = meta.Artist()
	}
	if meta.Album() != "" {
		s.Album = meta.Album()
	}
	if meta.Genre() != "" {
		s.Genre = meta.Genre()
	}
	if y := meta.Year(); y != 0 {
		s.Year = y
	}
	// rating/playcount - custom tags rarely available; leave as default
	return nil
}
