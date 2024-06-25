package snapshotbitcake

import (
	"Snapshot/app"
	"Snapshot/servent/message"
	"encoding/json"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type DelayedMessageSender struct {
	message message.Message
}

func (delayedMessageSender DelayedMessageSender) Run() {
	// We seed the random number generator here.
	rand.New(rand.NewSource(time.Now().Unix()))

	// Sleep for random duration before sending.
	time.Sleep(time.Duration(time.Duration(rand.Float64()*1000 + 500).Milliseconds()))

	// Used mostly for debugging.
	if MESSAGE_UTIL_PRINTING {
		switch delayedMessageSender.message.(type) {
		case *TransactionMessage:
			if msg, ok := delayedMessageSender.message.(*TransactionMessage); ok {
				app.TimeStampPrint("Sending message " + msg.ToString())
			}
		case *message.BasicMessage:
			if msg, ok := delayedMessageSender.message.(*message.BasicMessage); ok {
				app.TimeStampPrint("Sending message " + msg.ToString())
			}
		case *LYMarkerMessage:
			if msg, ok := delayedMessageSender.message.(*LYMarkerMessage); ok {
				app.TimeStampPrint("Sending message " + msg.ToString())
			}
		default:
			app.TimeStampErrorPrint("Message not recognized!")
			return
		}
	}

	app.ColorLock.Lock()
	defer app.ColorLock.Unlock()

	if !app.IsWhite.Load() {
		delayedMessageSender.message.SetRedColor()
	}

	var socket net.Conn
	var err error

	// So we don't repeat this useless casting.
	var receiverPort string
	var receiverIp string

	switch delayedMessageSender.message.(type) {
	case *TransactionMessage:
		if msg, ok := delayedMessageSender.message.(*TransactionMessage); ok {
			socket, err = net.Dial("tcp", msg.ReceiverInfo.Ip+":"+strconv.Itoa(msg.ReceiverInfo.Port))
			receiverPort = strconv.Itoa(msg.ReceiverInfo.Port)
			receiverIp = msg.ReceiverInfo.Ip
		}
	case *message.BasicMessage:
		if msg, ok := delayedMessageSender.message.(*message.BasicMessage); ok {
			socket, err = net.Dial("tcp", msg.ReceiverInfo.Ip+":"+strconv.Itoa(msg.ReceiverInfo.Port))
			receiverPort = strconv.Itoa(msg.ReceiverInfo.Port)
			receiverIp = msg.ReceiverInfo.Ip
		}
	case *LYMarkerMessage:
		if msg, ok := delayedMessageSender.message.(*LYMarkerMessage); ok {
			socket, err = net.Dial("tcp", msg.ReceiverInfo.Ip+":"+strconv.Itoa(msg.ReceiverInfo.Port))
			receiverPort = strconv.Itoa(msg.ReceiverInfo.Port)
			receiverIp = msg.ReceiverInfo.Ip
		}
	default:
		app.TimeStampErrorPrint("Message not recognized!")
		return
	}

	if err != nil {
		app.TimeStampErrorPrint("Failed to dialup " + receiverIp + " port: " + receiverPort)
		app.TimeStampErrorPrint(err.Error())
		return
	}

	msgToSend, errJson := json.Marshal(delayedMessageSender.message)

	if errJson != nil {
		app.TimeStampErrorPrint("Error converting message to json.")
		return
	}

	socket.Write(msgToSend)
	socket.Close()

	delayedMessageSender.message.SendEffect()

}
