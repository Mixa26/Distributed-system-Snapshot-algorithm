package handlersnapshot

import (
	"Snapshot/app"
	"Snapshot/app/snapshotbitcake"
	"Snapshot/servent/message"
)

type LYTellHandler struct {
	ClientMessage     message.Message
	SnapshotCollector snapshotbitcake.SnapshotCollector
}

func (lyTellHandler LYTellHandler) Run() {
	if msg, ok := lyTellHandler.ClientMessage.(*snapshotbitcake.LYTellMessage); ok {

		if msg.MessageType == "LY_TELL" {
			lyTellHandler.SnapshotCollector.AddLYSnapshotInfo(msg.OriginalSenderInfo.Id, msg.LYSnapshotResult)
		} else {
			app.TimeStampErrorPrint("Tell amount handler got: " + msg.MessageText)
		}

	} else {
		app.TimeStampErrorPrint("LYTellHandler received a Message that isn't of LYTellMessage type.")
	}
}
