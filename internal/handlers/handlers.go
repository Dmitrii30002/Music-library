package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/internal/models"
	DataBase "main/internal/storage"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// @Summary Add or get song information
// @Description Retrieves song information from the database. If the song is not found, it attempts to retrieve information from an external API and adds the song to the database.
// @Tags Songs
// @Accept  json
// @Produce  json
// @Param   group     query    string     true    "Group name"
// @Param   song      query    string     true    "Song title"
// @Success 200 {object} models.Song "Successfully found or added song"
// @Failure 400 {string} string "Bad request (missing 'group' or 'song')"
// @Failure 404 {string} string "Song not found in the database or external API"
// @Failure 406 {string} string "Failed to add new song to the database"
// @Router /songs [post]
func AddSong(c echo.Context) error {
	group := c.QueryParam("group")
	song := c.QueryParam("song")

	if song == "" || group == "" {
		log.Printf("[ERROR] Bad request")
		return c.String(http.StatusBadRequest, "Bad request")
	}

	db := DataBase.Connect()
	var helpSong models.Song
	err := db.Where("group = ? AND song = ?", group, song).First(&helpSong).Error
	if err != nil {
		log.Printf("[INFO] Song %s - %s not found. Try to get info from api", group, song)
		songWithDetails, finded := getSongDetailFromAPI(c, group, song)
		if !finded {
			return c.String(http.StatusNotFound, "Not found")
		}

		newSong := models.Song{
			Group:       group,
			Song:        song,
			ReleaseDate: songWithDetails.ReleaseDate,
			Text:        songWithDetails.Text,
			Link:        songWithDetails.Link,
		}

		err := db.Create(&newSong).Error
		if err != nil {
			log.Printf("[ERROR] Failed to add new song to DB: %v", err)
			return c.String(http.StatusNotAcceptable, "Not acceptable")
		}
		helpSong = newSong
		log.Printf("[INFO] Song %s - %s was added to DB", group, song)
	}
	return c.JSON(http.StatusOK, helpSong)
}

func getSongDetailFromAPI(c echo.Context, group string, song string) (models.SongDetail, bool) {
	encodedGroup := url.QueryEscape(group)
	encodedSong := url.QueryEscape(song)
	apiURL := fmt.Sprintf("http://localhost:8081/info?group=%s&song=%s", encodedGroup, encodedSong)
	response, err := http.Get(apiURL)
	if err != nil {
		log.Printf("[ERROR] Failed request to API: %v", err)
		c.String(http.StatusBadRequest, "Bad request")
		return models.SongDetail{}, false
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("[ERROR] Failed request to API: %v", err)
		c.String(http.StatusBadRequest, "Bad request")
		return models.SongDetail{}, true
	}

	var songDetail models.SongDetail
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("[ERROR] Failed read of responce: %v", err)
		c.String(http.StatusBadRequest, "Bad request")
		return models.SongDetail{}, true
	}

	if err := json.Unmarshal(body, &songDetail); err != nil {
		log.Printf("[ERROR] Failed parse of responce: %v", err)
		c.String(http.StatusBadRequest, "Bad request")
		return models.SongDetail{}, true
	}

	return songDetail, false
}

// @Summary Get a list of songs with pagination and filtering
// @Description Retrieves a list of songs from the database with support for pagination and filtering by group, song title, or release date.
// @Tags Songs
// @Accept  json
// @Produce  json
// @Param   group       query    string     false    "Filter by group name (optional)"
// @Param   song        query    string     false    "Filter by song title (optional)"
// @Param   releaseDate query    string     false    "Filter by release date (optional)"
// @Param   page        query    int        false    "Page number (default 1)" default(1)
// @Param   limit       query    int        false    "Number of items per page (default 10, max 100)" default(10)
// @Success 200 {array} models.Song "Successfully retrieved list of songs"
// @Failure 400 {string} string "Bad request (error parsing page or limit parameters)"
// @Router /songs [get]
func GetSongs(c echo.Context) error {
	group := c.QueryParam("group")
	song := c.QueryParam("song")
	releaseDate := c.QueryParam("releaseDate")

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		log.Printf("[ERROR] Wrong number of page: %v", err)
		page = 1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		log.Printf("[ERROR] Wrong number of limit: %v", err)
		limit = 10
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	db := DataBase.Connect()
	var songs []models.Song
	if group != "" {
		db.Order("group").Offset((page - 1) * limit).Limit(limit).Find(&songs)
	}
	if song != "" {
		db.Order("song").Offset((page - 1) * limit).Limit(limit).Find(&songs)
	}
	if releaseDate != "" {
		db.Order("release_date").Offset((page - 1) * limit).Limit(limit).Find(&songs)
	}
	if group == "" && song == "" && releaseDate == "" {
		db.Offset((page - 1) * limit).Limit(limit).Find(&songs)
	}

	log.Printf("[INFO] Song was given with pagination")
	return c.JSON(http.StatusOK, songs)
}

// @Summary Get song text with pagination by ID
// @Description Retrieves song text from the database by ID, splits it into verses, and returns the specified range of verses with pagination.
// @Tags Songs
// @Accept  json
// @Produce  json
// @Param   id      path    int     true    "Song ID"
// @Param   page    query   int     false   "Page number (default 1)" default(1)
// @Param   limit   query   int     false   "Number of verses per page (default 1)" default(1)
// @Success 200 {object} object  "Successfully retrieved song text"
// @Failure 400 {string} string "Bad request (error parsing page or limit parameters)"
// @Failure 404 {string} string "Song with the specified ID not found"
// @Router /songs/:id/lyrics [get]
func GetSongTextWithPagination(c echo.Context) error {
	id := c.Param("id")

	db := DataBase.Connect()

	var song models.Song
	if err := db.First(&song, id); err != nil {
		log.Printf("[ERROR] Song with id - %s was not founded", id)
		return c.String(http.StatusNotFound, "Not found")
	}

	verse := strings.Split(song.Text, "\n\n")

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		log.Printf("[ERROR] Wrong number of page: %v", err)
		page = 1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		log.Printf("[ERROR] Wrong number of limit: %v", err)
		limit = 1
	}

	start := (page - 1) * limit
	if start >= len(verse) {
		start = 0
	}
	end := start + limit
	if end > len(verse) {
		end = len(verse)
	}

	response := struct {
		Id     int64    `json:"ID"`
		Song   string   `json:"song"`
		Group  string   `json:"group"`
		Verses []string `json:"verse"`
		Total  int      `json:"total"`
	}{
		Id:     song.ID,
		Song:   song.Song,
		Group:  song.Group,
		Verses: verse[start:end],
		Total:  len(verse),
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary Update song information by ID
// @Description Updates existing song information in the database based on its ID.
// @Tags Songs
// @Accept  json
// @Produce  json
// @Param   id      path    int     true    "Song ID"
// @Param   song  body  models.Song true "Updated song information"
// @Success 200 {object} models.Song "Successfully updated song"
// @Failure 400 {string} string "Bad request (invalid song data)"
// @Failure 404 {string} string "Song with the specified ID not found"
// @Failure 500 {string} string "Internal Server Error (failed to update song)"
// @Router /songs/:id [put]
func UpdateSong(c echo.Context) error {
	id := c.Param("id")
	db := DataBase.Connect()

	err := db.First(id).Error
	if err != nil {
		log.Printf("[ERROR] Song with id %s was not found: %v", id, err)
		return c.String(http.StatusNotFound, "Not found")
	}

	var song models.Song
	err = c.Bind(&song)
	if err != nil {
		log.Printf("[ERROR] Invalid song: %v", err)
		return c.String(http.StatusBadRequest, "Bad request")
	}

	err = db.Save(&song).Error
	if err != nil {
		log.Printf("[ERROR] Song with id %s was not updated: %v", id, err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	log.Printf("[INFO] Song with id: %s was updated", id)
	return c.JSON(http.StatusOK, song)
}

// @Summary Delete song by ID
// @Description Deletes a song from the database based on its ID.
// @Tags Songs
// @Accept  json
// @Produce  plain
// @Param   id      path    int     true    "Song ID"
// @Success 200 {string} string "Successfully deleted song"
// @Failure 404 {string} string "Song was not found"
// @Router /songs/:id [delete]
func DeleteSong(c echo.Context) error {
	id := c.Param("id")
	db := DataBase.Connect()

	err := db.Delete(&models.Song{}, id).Error
	if err != nil {
		log.Printf("[ERROR] Failed to delete song with ID: %s %v", id, err)
		return c.String(http.StatusNotFound, "Song was not found")
	}

	log.Printf("[INFO] Song with id: %s was deleted", id)
	return c.String(http.StatusOK, "id #"+string(id)+": deleted")
}
