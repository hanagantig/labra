package firebase

import (
	"firebase.google.com/go/v4/messaging"
	"labra/internal/entity"
)

func NewMessageFromEntity(msg entity.Message) messaging.Message {
	return messaging.Message{
		Data:  msg.Data,
		Topic: msg.Topic,
	}
}
