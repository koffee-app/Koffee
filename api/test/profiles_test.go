package test

import (
	"fmt"
	"koffee/internal/config"
	"koffee/internal/models"
	"testing"
)

func BenchmarkGetProfilesImages(b *testing.B) {
	db := config.InitConfig()
	imgRepo := models.InitializeImages(db)
	repo := models.InitializeProfile(db, imgRepo)
	artists := repo.GetProfiles(&models.Profile{Artist: true}, true, 100)
	artistsNew := repo.GetProfilesImages(*artists)
	for _, artist := range artistsNew {
		fmt.Println(artist.ProfileImage, artist.UserID)
	}
}

func BenchmarkGetProfilesByIDs(b *testing.B) {
	db := config.InitConfig()
	imgRepo := models.InitializeImages(db)
	repo := models.InitializeProfile(db, imgRepo)
	artists := repo.GetProfilesByIDs([]uint32{2, 12821812, 2, 12821812, 2, 12821812, 3, 4, 5, 6, 5, 6, 5, 6, 5, 6, 5, 6, 5, 6, 5, 6, 5, 6, 5, 6, 5, 6, 5, 6, 5, 6, 5, 6, 5, 6, 2, 12821812})
	fmt.Println(artists)
}
