package snapshotbitcake

import (
	"Snapshot/app"
	"Snapshot/servent/message"
	"strconv"
)

type LYMarkerMessage struct {
	// This weird way of passing a basic message is for class extending(this class extends BasicMessage).
	message.BasicMessage
}

func ConstructLYMarkerMessage(
	sender *app.ServentInfo,
	receiver *app.ServentInfo,
	collectorId int,
) LYMarkerMessage {
	return LYMarkerMessage{
		// Pass to the "super" constructor.
		*message.ConstructBasicMessage("LY_MARKER", sender, receiver, strconv.Itoa(collectorId)),
	}
}
