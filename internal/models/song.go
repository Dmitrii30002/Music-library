package models

import (
	"time"

	"gorm.io/gorm"
)

type Song struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	Group       string    `json:"group"`
	Song        string    `json:"song"`
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

type SongDetail struct {
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

type SongStorage struct {
	db *gorm.DB
}

func NewSongStorage(db *gorm.DB) *SongStorage {
	return &SongStorage{db: db}
}

func (db *SongStorage) AddSong(song *Song) (*Song, error) {
	err := db.db.Create(song).Error
	if err != nil {
		return nil, err
	}

	return song, nil
}

func (db *SongStorage) GetAllSongs(page int, limit int) ([]Song, error) {
	var songs []Song
	offset := (page - 1) * limit
	err := db.db.Limit(limit).Offset(offset).Find(&songs).Error
	if err != nil {
		return nil, err
	}
	return songs, nil
}

func (db *SongStorage) GetSongByID(id int64) (*Song, error) {
	var song Song
	err := db.db.First(&song, id).Error
	if err != nil {
		return nil, err
	}
	return &song, nil
}

func (db *SongStorage) UpdateSong(song *Song) (*Song, error) {
	err := db.db.Updates(song).Error
	if err != nil {
		return nil, err
	}
	return song, nil
}

func (db *SongStorage) DeleteSongById(id int64) error {
	err := db.db.Delete(id).Error
	if err != nil {
		return err
	}
	return nil
}
