package models

import (
	"fmt"
	"koffee/pkg/formatter"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Image representation in the database
type Image struct {
	ID     uint32 `db:"id" json:"id"`
	XlURL  string `db:"xlurl" json:"xlurl"`
	MedURL string `db:"medurl" json:"medurl"`
	SmURL  string `db:"smurl" json:"smurl"`
	// ProfileImage, CoverImage, HeaderImage, HeaderImageAlbum...
	Type string `db:"type" json:"type"`
}

// ImageError .
type ImageError struct {
	Internal string `json:"internal"`
}

const schemaImg = `
	CREATE TABLE images (
		id integer,
		xlurl text NOT NULL,
		medurl text NOT NULL,
		smurl text NOT NULL,
		type text,
		UNIQUE (id, type)
	)
`

// ImageTypes is the available images
type ImageTypes uint8

// Images represents function for multiple images
type Images []Image

// Zip zips an array of images into a dictionary like { idOfImage: { ProfileImage: *Image } }
func (i Images) Zip() map[uint32]map[ImageTypes]Image {
	dict := make(map[uint32]map[ImageTypes]Image, len(i))
	for _, img := range i {
		mp, ok := dict[img.ID]
		if !ok {
			dict[img.ID] = map[ImageTypes]Image{StringToType(img.Type): img}
		} else {
			mp[StringToType(img.Type)] = img
		}
	}
	return dict
}

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

var stringToType = map[string]ImageTypes{"profile_image": ProfileImage, "cover_image": CoverImage, "header_image": HeaderImage, "header_image_album": HeaderImageAlbum}

// StringToType parses a string into a image type
func StringToType(s string) ImageTypes {
	return stringToType[s]
}

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
func InitializeImages(db *sqlx.DB) RepositoryImages {
	t := db.MustBegin()
	t.Exec(schemaImg)
	t.Commit()
	return &imageRepository{db: db}
}

// CreateOrUpdateImage tries to find an image in the table, if it exists it updates it
func (i *imageRepository) CreateOrUpdateImage(id uint32, urlXL, urlMed, urlSm string, typeImage ImageTypes) (*Image, *ImageError) {
	typeImg := getImageType(typeImage)
	if typeImg == "" {
		return nil, &ImageError{}
	}
	tx := i.db.MustBegin()
	_, err := tx.Exec("INSERT INTO images (id, xlurl, medurl, smurl, type) VALUES ($1, $2, $3, $4, $5) ON CONFLICT  (id, type) DO UPDATE SET xlurl=EXCLUDED.xlurl, medurl=EXCLUDED.medurl, smurl=$4", id, urlXL, urlMed, urlSm, typeImg)
	if err != nil {
		return nil, &ImageError{Internal: err.Error()}
	}
	err = tx.Commit()
	if err != nil {
		return nil, &ImageError{Internal: err.Error()}
	}
	return &Image{ID: id, XlURL: urlXL, MedURL: urlMed, SmURL: urlSm, Type: typeImg}, nil
}

func (i *imageRepository) DeleteImage(id uint32, typeImage ImageTypes) *ImageError {
	typeImg := getImageType(typeImage)
	if typeImg == "" {
		return &ImageError{}
	}
	tx := i.db.MustBegin()
	var uid uint32
	row := tx.QueryRowx("DELETE FROM images WHERE id=$1 AND type=$2 RETURNING id", id, typeImg).Scan(&uid)

	if row != nil && row.Error() != "" {
		return &ImageError{Internal: row.Error()}
	}

	if err := tx.Commit(); err != nil {
		return &ImageError{Internal: err.Error()}
	}

	if uid != id {
		return &ImageError{Internal: fmt.Sprintf("Uid is not equal to id, %d != %d", uid, id)}
	}

	return nil
}

func (i *imageRepository) GetImage(id uint32, typeImage ImageTypes) *Image {
	typeImg := getImageType(typeImage)
	if typeImg == "" {
		return nil
	}
	var images []Image
	tx := i.db.MustBegin()
	err := tx.Select(&images, "SELECT * FROM images WHERE id=$1 AND type=$2", id, typeImg)
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

func (i *imageRepository) GetImagesSameID(id uint32, typeImages ...ImageTypes) ([]Image, *ImageError) {
	typeImagesStr := make([]string, len(typeImages))
	for i := range typeImages {
		typeImagesStr[i] = getImageType(typeImages[i])
	}
	queryStr, array := formatter.Array(len(typeImages), "type", typeImagesStr)
	tx := i.db.MustBegin()
	var images []Image
	err := tx.Select(&images, fmt.Sprintf("SELECT * FROM images WHERE %s and id=%d", queryStr, id), array...)
	if err != nil {
		return nil, &ImageError{Internal: err.Error()}
	}
	err = tx.Commit()
	if err != nil {
		return nil, &ImageError{Internal: err.Error()}
	}
	return images, nil
}

// GetImagesFromIDs returns all the images specified in the ids and of the specified types
func (i *imageRepository) GetImagesFromIDs(typeImages []ImageTypes, ids ...uint32) ([]Image, error) {
	typeImagesStr := make([]string, len(typeImages))
	for i := range typeImages {
		typeImagesStr[i] = getImageType(typeImages[i])
	}
	queryStr, array := formatter.Array(len(typeImages), "type", typeImagesStr)
	tx := i.db.MustBegin()
	query := fmt.Sprintf("SELECT * FROM images WHERE id=ANY($%d) AND %s ORDER BY id ASC", len(array)+1, queryStr)
	array = append(array, pq.Array(ids))
	var images []Image
	err := tx.Select(&images, query, array...)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return images, nil
}
