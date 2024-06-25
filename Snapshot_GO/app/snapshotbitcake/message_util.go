package snapshotbitcake

import (
	"Snapshot/app"
	"Snapshot/servent/message"
	"encoding/json"
	"net"
)

// Helps with debugging.
var MESSAGE_UTIL_PRINTING bool = true

func ReadMessage(socket net.Conn) message.Message {

	buf := make([]byte, 1024)

	n, err := socket.Read(buf)
	if err != nil {
		app.TimeStampErrorPrint("Error in reading socket on: " + socket.LocalAddr().String())
	}

	var msg message.BasicMessage
	if err := json.Unmarshal(buf[:n], &msg); err != nil {
		app.TimeStampErrorPrint("Error reading json object from socker!")
	}

	if MESSAGE_UTIL_PRINTING {
		app.TimeStampPrint("Got message: " + msg.ToString())
	}

	return &msg
}

func SendMessage(msg message.Message) {
	DelayedMessageSender{msg}.Run()
}
