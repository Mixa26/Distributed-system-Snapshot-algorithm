package snapshotbitcake

import (
	"Snapshot/app"
	"Snapshot/servent/message"
	"strconv"
)

type TransactionMessage struct {
	// This weird way of passing a basic message is for class extending(this class extends BasicMessage).
	message.BasicMessage
	bitcakeManager BitcakeManager
}

func ConstructTransactionMessage(
	sender *app.ServentInfo,
	receiver *app.ServentInfo,
	amount int,
	bitcakeManager BitcakeManager,
) *TransactionMessage {
	newTransactionMessage := TransactionMessage{
		// Pass to the "super" constructor.
		BasicMessage:   *message.ConstructBasicMessage("TRANSACTION", sender, receiver, strconv.Itoa(amount)),
		bitcakeManager: bitcakeManager,
	}

	return &newTransactionMessage
}

/*
* We want to take away our amount exactly as we are sending, so our snapshots don't mess up.
* This method is invoked by the sender just before sending, and with a lock that guarantees
* that we are white when we are doing this in Chandy-Lamport.
 */
func (transactionMessage *TransactionMessage) SendEffect() {
	amount, err := strconv.Atoi(transactionMessage.MessageText)

	if err != nil {
		app.TimeStampErrorPrint("Passed a non number as amount int transaction message: " + transactionMessage.MessageText)
		return
	}

	transactionMessage.bitcakeManager.TakeSomeBitcakes(amount)
	if bm, ok := transactionMessage.bitcakeManager.(*LaiYangBitcakeManager); ok && transactionMessage.BasicMessage.White {
		bm.RecordGiveTransaction(transactionMessage.ReceiverInfo.Id, amount)
	}
}
