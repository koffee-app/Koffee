package controllers

import (
	"encoding/json"
	"fmt"
	"koffee/internal/middleware"
	"koffee/internal/models"
	"koffee/internal/rabbitmq"
	"koffee/internal/view"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/streadway/amqp"
)

type albumBody struct {
	Name        string `json:"name,omitempty"`
	CoverURL    string `json:"coverurl,omitempty"`
	Description string `json:"description,omitempty"`
	Published   bool   `json:"published,omitempty"`
	// *NOTE: Do we _really_ need this? We might just put this from our side and call it a day.
	Artists []string `json:"artists,omitempty"`
	// *WARNING: We might encounter integer overflows if we don't specify integer sth
	UserID        uint32   `json:"user_id,omitempty"`
	ID            uint32   `json:"id,omitempty"`
	AlbumID       uint32   `json:"album_id"`
	Collaborators []string `json:"collaborators,omitempty"`
}

func (a *albumBody) sanitize() {
	a.Name = strings.TrimSpace(a.Name)
	a.Description = strings.TrimSpace(a.Description)
	for i := range a.Artists {
		a.Artists[i] = strings.TrimSpace(a.Artists[i])
	}
}

type albumController struct {
	repository  models.AlbumsRepository
	mq          rabbitmq.MessageListener
	albumSender rabbitmq.MessageSender
}

// InitializeAlbumController inits the controller of albums
func InitializeAlbumController(routes *Group, router *httprouter.Router, repo models.AlbumsRepository, q rabbitmq.MessageListener) {
	albumImpl := albumController{repository: repo, mq: q, albumSender: q.NewSender("new_album")}

	q.OnMessage("update_collaborators", albumImpl.updateCollab)

	group := New(routes, "/albums")

	// Gets every owned album by the user
	router.GET(group.Route("/owned"), middleware.Authorization(albumImpl.getOwnedAlbums))
	// Get a owned album by the user
	router.GET(group.Route("/owned/:album_id"), middleware.Authorization(albumImpl.getOwnedAlbum))
	router.GET(group.Route("/public/:profile"), albumImpl.getPublicAlbums)
	router.GET(group.Route("/id/:album_id"), albumImpl.getPublicAlbum)
	router.PUT(group.Route("/:id"), middleware.Authorization(albumImpl.updateAlbum))
	router.GET(group.Route("/full/:id"), albumImpl.getAlbumFull)

	// Post an album
	router.POST(group.Route(""), middleware.Authorization(albumImpl.createAlbum))

	// Get a public album

}

func (a *albumController) getAlbumFull(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	IDu64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		view.AlbumError(w, &models.AlbumError{Internal: "ID is not an integer"})
		return
	}
	album, err := a.repository.GetAlbumFull(uint32(IDu64), true)
	if err != nil {
		view.AlbumError(w, &models.AlbumError{Internal: err.Error()})
		return
	}
	view.PrepareJSON(w, 200)
	view.ReturnJSON(w, view.JSON(album, "Success", 200))
}

func (a *albumController) createAlbum(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := middleware.GetUser(r)
	body := albumBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		view.AlbumError(w, &models.AlbumError{Internal: "Bad body request"})
		return
	}
	body.sanitize()
	if album, albumErr := a.repository.CreateAlbum(user.UserID, body.Name, body.Artists, body.Description); albumErr != nil {
		view.AlbumError(w, albumErr)
	} else {
		view.Album(w, album)
		album.EmailCreator = user.Email
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
	var body albumBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		view.Error(w, "Error parsing request body", http.StatusBadRequest, err)
		return
	}
	body.sanitize()
	publish := strconv.FormatBool(body.Published)
	albumID := params.ByName("id")
	IDu64, _ := strconv.ParseUint(albumID, 10, 32)
	album, err := a.repository.UpdateAlbum(user.UserID, uint32(IDu64), publish, body.Description, body.Name, "")
	if err != nil {
		view.AlbumError(w, err)
		return
	}
	view.Album(w, album)
}

func (a *albumController) updateCollab(msg *amqp.Delivery) {
	var body albumBody
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		log.Println(err)
		return
	}

	a.repository.NewCollaborators(body.Artists, body.ID)
}

func (a *albumController) changeCoverURL(msg *amqp.Delivery) {
	var msgBody albumBody
	if err := json.Unmarshal(msg.Body, &msgBody); err == nil {
		album := a.repository.NewCollaborators(msgBody.Collaborators, msgBody.ID)
		log.Println("Update collab in album ", album.Artists)
	}
}

// (after, beforeN, nItems)
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
