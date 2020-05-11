package controllers

import (
	"encoding/json"
	"koffee/internal/models"
	"koffee/internal/rabbitmq"
	"log"

	"github.com/streadway/amqp"
)

type imageController struct {
	listener rabbitmq.MessageListener
	repo     models.ImagesRepository
}

type bodyImage struct {
	// We use uints between the queues because we will inform in the docs what type is each number and better performance
	Type   models.ImageTypes `json:"image_type"`
	URLxl  string            `json:"url_xl"`
	URLmed string            `json:"url_med"`
	URLsm  string            `json:"url_sm"`
	ID     uint32            `json:"id"`
}

// InitializeImageController inits the rabbit controller for receiving all of the messages from the image services;
// this means that every change to the image urls or if there is a new image in the image service it will be notified here
func InitializeImageController(l rabbitmq.MessageListener, repo models.ImagesRepository) {
	ic := imageController{listener: l, repo: repo}
	l.OnMessage("new_image_update", ic.handler)
}

func (i *imageController) handler(m *amqp.Delivery) {
	var target bodyImage
	if err := json.Unmarshal(m.Body, &target); err != nil {
		return
	}
	img, err := i.repo.CreateOrUpdateImage(target.ID, target.URLxl, target.URLmed, target.URLsm, target.Type)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	log.Printf("Success: %v\n", img)
}
