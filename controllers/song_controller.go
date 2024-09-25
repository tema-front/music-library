package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"music-library/models"
	"music-library/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

var db *gorm.DB

const (
	defaultLimit = 10
	defaultOffset = 0
)

func SetDB(database *gorm.DB) {
	db = database
}

func GetSongs(w http.ResponseWriter, r *http.Request) {
	var songs []models.Song

	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")

	limit, offset := utils.ParseLimitAndOffset(r, defaultLimit, defaultOffset)

	query := db.Model(&models.Song{})
	if group != "" {
    query = query.Where("\"group\" ILIKE ?", "%"+group+"%")
	}
	if song != "" {
		query = query.Where("\"song\" ILIKE ?", "%"+song+"%")
	}

	if err := query.Offset(offset).Limit(limit).Find(&songs).Error; err != nil {
		http.Error(w, fmt.Sprintf("couldn't get songs: %v", err), http.StatusInternalServerError)
		return
	}

	response := songs
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetSongText(w http.ResponseWriter, r *http.Request) {
	songIDStr := chi.URLParam(r, "id")
	songID, err := strconv.Atoi(songIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't parse songID: %v", err), http.StatusBadRequest)
		return
	}

	var song models.Song
	if result := db.First(&song, songID); result.Error != nil {
		http.Error(w, fmt.Sprintf("couldn't get song: %v", result.Error), http.StatusInternalServerError)
		return
	}

	response := []string{}
	if len(song.Text) == 0 {
		json.NewEncoder(w).Encode(response)
		return
	}

	verses := strings.Split(song.Text, "\n\n")
	limit, offset := utils.ParseLimitAndOffset(r, defaultLimit, defaultOffset)
	
	response = verses
	response, err = utils.SafeSlice(offset, limit, verses)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't parse limit and offset: %v", err), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteSong(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := db.Delete(&models.Song{}, id).Error; err != nil {
		http.Error(w, fmt.Sprintf("couldn't delete song: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func EditSong(w http.ResponseWriter, r *http.Request) {
	songIDStr := chi.URLParam(r, "id")
	songID, err := strconv.Atoi(songIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't parse songID: %v", err), http.StatusBadRequest)
		return
	}

	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, fmt.Sprintf("couldn't parse song: %v", err), http.StatusBadRequest)
		return
	}

	song.ID = uint(songID)
	if err := db.Save(&song).Error; err != nil {
		http.Error(w, fmt.Sprintf("couldn't update song: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(song)
}

func AddSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, fmt.Sprintf("crror parsing JSON: %v", err), http.StatusBadRequest)
		return
	}

	info, err := getSongDetail(song.Group, song.Song)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't get song detail: %v", err), http.StatusInternalServerError)
		return
	}

	song.ReleaseDate, song.Text, song.Link = info.ReleaseDate, info.Text, info.Link

	if err := db.Create(&song).Error; err != nil {
		http.Error(w, fmt.Sprintf("couldn't create song: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getSongDetail(group, song string) (*models.SongDetail, error) {
	url := fmt.Sprintf("http://api.example.com/info?group=%s&song=%s", group, song)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	var songDetail models.SongDetail
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &songDetail); err != nil {
		return nil, err
	}

	return &songDetail, nil
}