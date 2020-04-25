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
		b.Errorf("Image not equal to the created or non existing")
	}
	fmt.Println(imgGet, img)
	// deletion := repo.DeleteImage(img.ID, models.CoverImage)
	// if deletion != nil {
	// 	fmt.Println(deletion)
	// 	panic("REMEMBER to delete this image because there was a mistake when deleting")
	// }
}

func BenchmarkIn(b *testing.B) {
	db := config.InitConfig()
	models.InitializeImages(db).GetImagesFromIDs([]models.ImageTypes{models.CoverImage, models.ProfileImage}, 1, 2, 3, 4)
}
func BenchmarkImagesCreation(b *testing.B) {
	db := config.InitConfig()
	repo := models.InitializeImages(db)
	repoProfiles := models.InitializeProfile(db, repo)

	img, err := repo.CreateOrUpdateImage(2, "longurl", "mediumurl", "smallurl", models.CoverImage)
	imgProfile, err := repo.CreateOrUpdateImage(2, "longurlprofile2", "mediumurlprofile", "smallurlprofile", models.ProfileImage)
	if err != nil {
		panic(fmt.Errorf("Error: %#v", err))
	}
	imgs, err := repo.GetImagesSameID(2, models.CoverImage, models.ProfileImage)
	if err != nil {
		panic(err)
	}
	if len(imgs) < 2 {
		panic(fmt.Errorf("Error, image length is too short: %d", len(imgs)))
	}
	if imgs[0].ID == img.ID && img.XlURL == imgs[0].XlURL && img.Type == imgs[0].Type {
		b.Log("Correct first image")
	} else {
		panic("Incorrect first image")
	}
	if imgs[1].ID == imgProfile.ID && imgProfile.XlURL == imgs[1].XlURL && imgProfile.Type == imgs[1].Type {
		b.Log("Correct second image")
	} else {
		panic("Incorrect second image")
	}
	b.Log(imgs)
	profile := repoProfiles.GetImage(&models.Profile{UserID: 2})
	if profile.ProfileImage != nil && profile.HeaderImage == nil {
		b.Log(profile.ProfileImage, profile.HeaderImage)
	} else {
		b.Log(profile.ProfileImage, profile.HeaderImage)
		panic("ProfileImage should be filled and HeaderImage should be empty")
	}
	// repo.DeleteImage(2, models.ProfileImage)
	// repo.DeleteImage(2, models.CoverImage)
}
