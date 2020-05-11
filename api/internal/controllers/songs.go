package controllers

import (
	"koffee/internal/models"
	"koffee/internal/view"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type songsBody struct {
}

type songsImpl struct {
	repository models.SongsRepository
}

// InitializeSongsController initializes the router of the songs controller
func InitializeSongsController(api *Group, router *httprouter.Router, repo models.SongsRepository) {
	s := songsImpl{repository: repo}
	group := New(api, "/songs")
	router.GET(group.Route("/:albumID"), s.GetSongsFromAlbum)
}

func (s *songsImpl) GetSongsFromAlbum(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	albumIDStr := p.ByName("albumID")
	albumIDUint64, err := strconv.ParseUint(albumIDStr, 10, 32)
	if err != nil {
		view.Error(w, "Error parsing the parameter albumID, not a string", 400, err)
		return
	}
	_, err = s.repository.GetSongsByID(uint32(albumIDUint64))
	if err == models.ErrSongNotFound {

	}
}
