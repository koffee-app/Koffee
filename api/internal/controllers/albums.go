package controllers

import (
	"encoding/json"
	"fmt"
	"koffee/internal/middleware"
	"koffee/internal/models"
	"koffee/internal/rabbitmq"
	"koffee/internal/view"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/streadway/amqp"
)

type albumBody struct {
	Name        string   `json:"name,omitempty"`
	CoverURL    string   `json:"coverurl,omitempty"`
	Description string   `json:"description,omitempty"`
	Published   bool     `json:"published,omitempty"`
	Artists     []string `json:"artists,omitempty"`
	UserID      uint32   `json:"user_id,omitempty"`
	ID          uint32   `json:"id,omitempty"`
}

type albumController struct {
	repository  models.RepositoryAlbums
	mq          rabbitmq.MessageListener
	albumSender rabbitmq.MessageSender
}

// InitializeAlbumController inits the controller of albums
func InitializeAlbumController(routes *Group, router *httprouter.Router, repo models.RepositoryAlbums, q rabbitmq.MessageListener) {
	albumImpl := albumController{repository: repo, mq: q, albumSender: q.NewSender("new_album")}

	q.OnMessage("new_cover_url", albumImpl.changeCoverURL)

	group := New(routes, "/albums")

	// Gets every owned album by the user
	router.GET(group.Route("/owned"), middleware.Authorization(albumImpl.getOwnedAlbums))
	// Get a owned album by the user
	router.GET(group.Route("/owned/:album_id"), middleware.Authorization(albumImpl.getOwnedAlbum))
	router.GET(group.Route("/public/:profile"), albumImpl.getPublicAlbums)
	router.GET(group.Route("/id/:album_id"), albumImpl.getPublicAlbum)
	router.PUT(group.Route("/:id"), middleware.Authorization(albumImpl.updateAlbum))

	// Post an album
	router.POST(group.Route(""), middleware.Authorization(albumImpl.createAlbum))

	// Get a public album

}

func (a *albumController) createAlbum(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := middleware.GetUser(r)
	body := albumBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		view.AlbumError(w, &models.AlbumError{Internal: "Bad body request"})
		return
	}
	if album, albumErr := a.repository.CreateAlbum(user.UserID, body.Name, body.Artists, body.Description); albumErr != nil {
		view.AlbumError(w, albumErr)
	} else {
		view.Album(w, album)
		// Send it to the listeners of this queue
		view.SendJSON(a.albumSender, album)
	}
}

func (a *albumController) getPublicAlbum(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	album := params.ByName("album_id")
	albumID, err := strconv.ParseUint(album, 10, 32)
	if err != nil {
		view.Error(w, "Error parsing id", 400, err)
		return
	}
	albumModel, errAlbum := a.repository.GetPublicAlbumByID(uint32(albumID))
	if errAlbum != nil {
		view.AlbumError(w, errAlbum)
		return
	}
	view.Album(w, albumModel)
}

func (a *albumController) getPublicAlbums(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	values := r.URL.Query()

	profile := params.ByName("profile")
	after, before, nItems := retrievePaginationValues(&values)

	albums, err := a.repository.GetProfileAlbumsByUsername(profile, true, after, before, nItems)

	if err != nil {
		view.AlbumError(w, err)
		return
	}
	view.Albums(w, albums)
}

func (a *albumController) getOwnedAlbums(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := middleware.GetUser(r)
	values := r.URL.Query()
	published := values.Get("published")
	after, before, nItems := retrievePaginationValues(&values)
	var albums []models.Album
	var err *models.AlbumError
	if published != "" {
		albums, err = a.repository.GetAlbumsByUserIDPublish(user.UserID, published == "true", after, before, nItems)

	} else {
		albums, err = a.repository.GetAlbumsByUserID(user.UserID, after, before, nItems)
	}
	if err == nil {
		view.Albums(w, albums)
	} else {
		view.AlbumError(w, err)
	}
}

func (a *albumController) getOwnedAlbum(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idStr := p.ByName("album_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		view.Error(w, "Error parsing ID", 400, err)
		return
	}
	album, errAlbum := a.repository.GetAlbumOwnedByID(uint32(id), middleware.GetUser(r).UserID)
	if errAlbum != nil {
		view.AlbumError(w, errAlbum)
		return
	}
	view.Album(w, album)
}

func (a *albumController) updateAlbum(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user := middleware.GetUser(r)
	values := r.URL.Query()
	des := strings.Trim(values.Get("description"), " ")
	publish := strings.Trim(values.Get("published"), " ")
	name := strings.Trim(values.Get("name"), " ")
	albumID := params.ByName("id")
	IDu64, _ := strconv.ParseUint(albumID, 10, 32)
	album, err := a.repository.UpdateAlbum(user.UserID, uint32(IDu64), publish, des, name, "")
	if err != nil {
		view.AlbumError(w, err)
		return
	}
	view.Album(w, album)
}

func (a *albumController) changeCoverURL(msg *amqp.Delivery) {
	var msgBody albumBody
	if err := json.Unmarshal(msg.Body, &msgBody); err == nil {
		a.repository.UpdateAlbum(msgBody.UserID, msgBody.ID, "", "", "", msgBody.CoverURL)
	}
}

// (after, nItems)
func retrievePaginationValues(values *url.Values) (int, int, int) {
	after := values.Get("after")
	nItems := values.Get("nitems")
	before := values.Get("before")

	afterN, _ := strconv.Atoi(after)
	nItemsN, _ := strconv.Atoi(nItems)
	beforeN, _ := strconv.Atoi(before)

	fmt.Println(afterN, nItemsN)
	return afterN, beforeN, nItemsN
}
