package handler

import (
	"Snapshot/app"
	"Snapshot/servent/message"
)

type NullHandler struct {
	ClientMessage message.Message
}

func (nullHandler NullHandler) Run() {
	if message, ok := nullHandler.ClientMessage.(*message.BasicMessage); ok {
		app.TimeStampErrorPrint("Couldn't handle message: " + message.MessageText)
	} else {
		app.TimeStampErrorPrint("Message wasn't a basic message!")
	}
}
