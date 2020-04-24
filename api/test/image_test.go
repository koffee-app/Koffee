package test

import (
	"fmt"
	"koffee/internal/config"
	"koffee/internal/models"
	"testing"
)


func BenchmarkImageCreation(b *testing.B) {
	db := config.InitConfig()
	repo := models.InitializeImages(db)
	img, err := repo.CreateOrUpdateImage(1, "longurl", "mediumurl", "smallurl", models.CoverImage)
	if err != nil {
		panic(fmt.Errorf("Error: %#v", err))
	}
	imgGet := repo.GetImage(img.ID, models.CoverImage)
	if imgGet == nil || imgGet.ID != img.ID || imgGet.Type != img.Type || imgGet.XlURL != img.XlURL {
		fmt.Println(imgGet, img)
		panic("Image not equal to the created or non existing")
	}
	fmt.Println(imgGet)
	deletion := repo.DeleteImage(img.ID, models.CoverImage)
	if deletion != nil {
		fmt.Println(deletion)
		panic("REMEMBER to delete this image because there was a mistake when deleting")
	}
}	