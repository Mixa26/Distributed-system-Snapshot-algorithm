package message

import (
	"Snapshot/app"
	"strconv"
	"sync/atomic"
)

type BasicMessage struct {
	MessageType        string
	OriginalSenderInfo *app.ServentInfo
	ReceiverInfo       *app.ServentInfo
	RouteList          []*app.ServentInfo
	MessageText        string
	White              bool
	MessageCounter     atomic.Int64
	MessageId          int
}

// Basic message constructor
func ConstructBasicMessage(MessageType string,
	OriginalSenderInfo *app.ServentInfo,
	ReceiverInfo *app.ServentInfo,
	MessageText string) *BasicMessage {
	newBasicMessage := BasicMessage{
		MessageType:        MessageType,
		OriginalSenderInfo: OriginalSenderInfo,
		ReceiverInfo:       ReceiverInfo,
		White:              app.IsWhite.CompareAndSwap(false, true),
		MessageText:        MessageText,
	}

	// Empty the slice if it has something.
	newBasicMessage.RouteList = nil

	newBasicMessage.MessageCounter.Add(1)

	return &newBasicMessage
}

// Adds us to the route list, doesn't change the original owner.
// Used for resending the message to the next receiver.
func (basicMessage *BasicMessage) MakeMeASender() Message {
	// Make a copy of the route list so we don't change
	// the existing one.
	var newRouteList []*app.ServentInfo
	copy(newRouteList, basicMessage.RouteList)

	newRouteList = append(newRouteList, &app.MyServentInfo)

	newMessage := basicMessage
	newMessage.RouteList = newRouteList

	return newMessage
}

/*
* Change the message received based on ID. The receiver has to be our neighbor.
* Use this when you want to send a message to multiple neighbors, or when resending.
 */
func (basicMessage *BasicMessage) ChangeReceiver(newReceiverId int) Message {
	// If we're trying to send a message to a non-neighboor, abort.
	for _, neighbor := range app.MyServentInfo.Neighbors {
		if neighbor == newReceiverId {
			// Set the new receiver.
			newMessage := basicMessage
			newMessage.ReceiverInfo = app.GetServentInfo(newReceiverId)

			return newMessage
		}
	}

	app.TimeStampErrorPrint("Trying to make a message for " + strconv.Itoa(newReceiverId) + " who is not a neighbor.")
	return nil
}

func (basicMessage *BasicMessage) SetRedColor() Message {
	newMessage := basicMessage
	newMessage.White = false

	return newMessage
}

func (basicMessage *BasicMessage) SetWhiteColor() Message {
	newMessage := basicMessage
	newMessage.White = true

	return newMessage
}

func (basicMessage *BasicMessage) Equals(basicMessage1 *BasicMessage) bool {
	return basicMessage.MessageId == basicMessage1.MessageId &&
		basicMessage.OriginalSenderInfo.Id == basicMessage1.OriginalSenderInfo.Id
}

func (basicMessage *BasicMessage) Hash() int {
	hash, _ := strconv.Atoi(strconv.Itoa(basicMessage.MessageId) + strconv.Itoa(basicMessage.OriginalSenderInfo.Id))
	return hash
}

func (basicMessage *BasicMessage) ToString() string {
	return "[" +
		strconv.Itoa(basicMessage.OriginalSenderInfo.Id) + " | " +
		strconv.Itoa(basicMessage.MessageId) + " | " +
		basicMessage.MessageText + " | " +
		basicMessage.MessageType + " | " +
		strconv.Itoa(basicMessage.ReceiverInfo.Id) +
		"]"
}

func (basicMessage *BasicMessage) SendEffect() {
	// Do nothing
}
