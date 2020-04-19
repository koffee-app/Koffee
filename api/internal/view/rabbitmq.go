package view

import (
	"encoding/json"
	"koffee/internal/rabbitmq"
)

// SendJSON accross a rabbitmq Queue
func SendJSON(m rabbitmq.MessageSender, i interface{}) {
	b, _ := json.Marshal(&i)
	m.Send(b, "application/json")
}
