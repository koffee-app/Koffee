package controllers

import (
	"encoding/json"
	"koffee/internal/models"
	"koffee/internal/rabbitmq"

	"github.com/streadway/amqp"
)

type imageController struct {
	listener rabbitmq.MessageListener
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
func InitializeImageController(l rabbitmq.MessageListener) {
	ic := imageController{listener: l}
	l.OnMessage("new_image_update", ic.handler)
}

func (i *imageController) handler(m *amqp.Delivery) {
	var target bodyImage

	if err := json.Unmarshal(m.Body, &target); err != nil {
		return
	}

	// We can use an array of handlers { func1(), func2() ... } but i think this would be a bit faster
	switch target.Type {
	case models.CoverImage:
		{
			break
		}
	case models.HeaderImage:
		{
			break

		}

	case models.HeaderImageAlbum:
		{
			break
		}

	case models.ProfileImage:
		{
			break
		}
	}
}
