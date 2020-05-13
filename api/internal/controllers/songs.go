package controllers

import (
	"koffee/internal/models"
	"koffee/internal/view"
	"log"
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
	router.GET(group.Route("/:id"), s.GetSong)
	router.GET(group.Route("/album/:albumID"), s.GetSongsFromAlbum)
}

// TODO Not tested
func (s *songsImpl) GetSong(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	songIDStr := p.ByName("albumID")
	songIDUint64, err := strconv.ParseUint(songIDStr, 10, 32)
	if err != nil {
		view.Error(w, "Error parsing the parameter id, not a number", 400, err)
		return
	}
	song, err := s.repository.GetSongByID(uint32(songIDUint64))
	if err == models.ErrSongNotFound {
		view.Error(w, "Not found", http.StatusNotFound, err)
		return
	}
	view.Song(w, song)
}

// TODO Not tested
func (s *songsImpl) GetSongsFromAlbum(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	albumIDStr := p.ByName("albumID")
	albumIDUint64, err := strconv.ParseUint(albumIDStr, 10, 32)
	if err != nil {
		view.Error(w, "Error parsing the parameter albumID, not a number", 400, err)
		return
	}
	songs, err := s.repository.GetSongsByID(uint32(albumIDUint64))

	if err == models.ErrSongNotFound {
		view.Error(w, "Error", 404, err)
		return
	}

	if err != nil {
		log.Println(err)
		view.Error(w, "Internal error", http.StatusBadRequest, err)
		return
	}

	view.Songs(w, songs)
}
