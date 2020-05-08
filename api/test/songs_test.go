package test

import (
	"koffee/internal/config"
	"koffee/internal/models"
	"log"
	"testing"
)

func BenchmarkGetSongs(b *testing.B) {
	db := config.InitConfig()
	repo := models.InitializeSongRepository(db)
	songs, err := repo.GetSongsByID(82)
	if err != nil {
		panic(err)
	}
	log.Println(songs)
}
