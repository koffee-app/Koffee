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
