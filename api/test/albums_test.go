package test

import (
	"fmt"
	"koffee/internal/config"
	"koffee/internal/models"
	"log"
	"testing"
)

func BenchmarkGetProfileAlbums(b *testing.B) {
	db := config.InitConfig()
	repo := models.InitializeAlbums(db, models.InitializeImages(db))
	albums, err := repo.GetProfileAlbumsByUsername("gabivlj", true, 0, 0, 100)
	if err != nil {
		panic(err)
	}
	log.Println(albums)
}

func BenchmarkGetAlbumsImages(b *testing.B) {
	db := config.InitConfig()

	imgRepo := models.InitializeImages(db)
	repo := models.InitializeAlbums(db, imgRepo)

	// // a, err := repo.CreateAlbum(2, "Incredible album", []string{"gabivlj"}, "An incredible album.")
	// if err != nil {
	// 	panic(err)
	// }
	img, errImg := imgRepo.CreateOrUpdateImage(63, "xddd", "xdd", "xd", models.CoverImage)
	fmt.Println(img.XlURL)
	defer func() {

		if errImg == nil {
			imgRepo.DeleteImage(img.ID, models.CoverImage)
		}
	}()
	if errImg != nil {
		panic(errImg)
	}

	albums, _ := repo.GetAlbumsByUserID(2, 0, 0, 120)

	albums = repo.GetAlbumsImages(albums)
	for _, album := range albums {
		fmt.Println(album.CoverImage, album.ID, album.Name)
	}

}
