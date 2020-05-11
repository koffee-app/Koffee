package repository

import (
	"fmt"
	"koffee/internal/models"
	"koffee/pkg/logger"

	"github.com/jmoiron/sqlx"
)

// TestCreateAlbum mini-test to check that everything works
// TODO Delete the db sqlx.DB ref
func TestCreateAlbum(p models.ProfilesRepository, ra models.AlbumsRepository, db *sqlx.DB) bool {
	defer func() {
		tx := db.MustBegin()
		_ = tx.QueryRowx("DELETE FROM profiles WHERE userid=23231323")
		_ = tx.QueryRowx("DELETE FROM profiles WHERE userid=23231324")
		tx.Commit()
	}()

	logger.Log("test_create", "Creating profile...")

	// We really don't care if this was succesful
	_, _ = p.CreateProfile("the_weeknd", "the_weeknd", 23231323, true)
	_, _ = p.CreateProfile("gabivlj2", "gabivlj2", 23231324, true)

	logger.Log("test_create", "Creating Album...")

	a, e := ra.CreateAlbum(23231324, "After Hours", []string{"the_weeknd", "gabivlj2"}, "After hours is a great album!")

	if e != nil {
		fmt.Println("ERROR CREATING ALBUM: ", *e)
		logger.Log("test_create", "Test did not pass!")
		return false
	}

	tx := db.MustBegin()

	tx.Exec("DELETE FROM albums WHERE id=$1", a.ID)
	logger.Log("test_create", tx.Commit())
	logger.Log("test_create", "Test passed 🥳")
	return a != nil && e == nil
}
