package test

import (
	"koffee/internal/config"
	"koffee/internal/models"
	"log"
	"testing"
)

func BenchmarkGetProfileAlbums(b *testing.B) {
	db := config.InitConfig()
	repo := models.InitializeAlbums(db)
	albums, err := repo.GetProfileAlbumsByUsername("gabivlj", true, 0, 0, 100)
	if err != nil {
		panic(err)
	}
	log.Println(albums)
}
