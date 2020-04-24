package models

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// Image representation in the database
type Image struct {
	ID     uint32 `db:"id"`
	XlURL  string `db:"xlurl"`
	MedURL string `db:"mdurl"`
	SmURL  string `db:"smurl"`
	// ProfileImage, CoverImage, HeaderImage, HeaderImageAlbum...
	Type string `db:"type"`
}

// ImageError .
type ImageError struct {
	Internal string `json:"internal"`
}

const schemaImg = `
	CREATE TABLE images (
		id integer,
		xlurl text,
		medurl text,
		smurl text,
		type text
	)
`

// ImageTypes is the available images
type ImageTypes uint8

const (
	// ProfileImage type
	ProfileImage ImageTypes = iota
	// CoverImage type
	CoverImage
	// HeaderImage type
	HeaderImage
	// HeaderImageAlbum type
	HeaderImageAlbum
)

var typeToString = []string{"profile_image", "cover_image", "header_image", "header_image_album"}

func getImageType(typeImg ImageTypes) string {
	if int(typeImg) >= len(typeToString) {
		return ""
	}
	return typeToString[typeImg]
}

type imageRepository struct {
	db *sqlx.DB
}

// InitializeImages .
func InitializeImages(db *sqlx.DB) {
	t := db.MustBegin()
	t.Exec(schemaImg)
	t.Commit()
	// return &imageRepository{db: db}
}

// CreateOrUpdateImage tries to find an image in the table, if it exists it updates it
func (i *imageRepository) CreateOrUpdateImage(id uint32, urlXL, urlMed, urlSm string, typeImage ImageTypes) (*Image, *ImageError) {
	typeImg := getImageType(typeImage)
	if typeImg == "" {
		return nil, &ImageError{}
	}
	tx := i.db.MustBegin()
	_, err := tx.Exec("INSERT INTO images (id, xlurl, medurl, smurl, type) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id, type) DO UPDATE SET xlurl=$2, medurl=$3, smurl=$4", id, urlXL, urlMed, urlSm, typeImg)
	if err != nil {
		return nil, &ImageError{Internal: err.Error()}
	}
	err = tx.Commit()
	if err != nil {
		return nil, &ImageError{Internal: err.Error()}
	}
	return &Image{ID: id, XlURL: urlXL, MedURL: urlMed, SmURL: urlSm, Type: typeImg}, nil
}

func (i *imageRepository) GetImage(id uint32, typeImage ImageTypes) *Image {
	typeImg := getImageType(typeImage)
	if typeImg == "" {
		return nil
	}
	var images []Image
	tx := i.db.MustBegin()
	err := tx.Select(&images, "SELECT * FROM images WHERE id=$1 AND type=$2")
	if err != nil {
		log.Println(err)
		return nil
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return nil
	}
	if len(images) == 0 {
		return nil
	}
	return &images[0]
}
