package main

import (
	"log"
	_ "main/docs"
	"main/internal/handlers"
	"main/internal/models"
	DataBase "main/internal/storage"

	"github.com/joho/godotenv"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
)

// @title Music Library API
// @version 1.0
// @description API of music library
// @host localhost:8080
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("[ERROR] Failed to load config: %v", err)
	}

	db := DataBase.Connect()
	log.Println("[INFO] Database connection established")

	db.AutoMigrate(&models.Song{})
	log.Println("[INFO] Database was migrations completed")

	e := echo.New()

	e.GET("/songs", handlers.GetSongs)
	e.GET("/songs/:id/lyrics", handlers.GetSongTextWithPagination)
	e.DELETE("/songs/:id", handlers.DeleteSong)
	e.PUT("/songs/:id", handlers.UpdateSong)
	e.POST("/songs", handlers.AddSong)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	log.Println("[INFO] Server was started at 8080!!!")
	e.Start(":8080")

}
