package snapshotbitcake

import (
	"Snapshot/app"
	"Snapshot/servent/message"
)

type LYTellMessage struct {
	// This weird way of passing a basic message is for class extending(this class extends BasicMessage).
	message.BasicMessage
	LYSnapshotResult *LYSnapshotResult
}

// Constructor for LYTellMessage.
func ConstructLYTellMessage(
	sender *app.ServentInfo,
	receiver *app.ServentInfo,
	lySnapshotResult *LYSnapshotResult) *LYTellMessage {
	newLYTellMessage := LYTellMessage{
		// Pass to the "super" constructor.
		BasicMessage:     *message.ConstructBasicMessage("LY_TELL", sender, receiver, ""),
		LYSnapshotResult: lySnapshotResult,
	}

	return &newLYTellMessage
}

func (lyTellMessage *LYTellMessage) SetRedColor() message.Message {
	newLYTellMessage := lyTellMessage
	newLYTellMessage.White = false

	return newLYTellMessage
}
