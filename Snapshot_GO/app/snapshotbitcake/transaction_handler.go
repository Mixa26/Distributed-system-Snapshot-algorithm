package snapshotbitcake

import (
	"Snapshot/app"
	"Snapshot/servent/message"
	"strconv"
)

type TransactionalHandler struct {
	ClientMessage  message.Message
	BitcakeManager BitcakeManager
}

func (transactionalHandler TransactionalHandler) Run() {
	if msg, ok := transactionalHandler.ClientMessage.(*message.BasicMessage); ok {
		if msg.MessageType == "TRANSACTION" {
			// Fetch the amount used in the transaction.
			amount, err := strconv.Atoi(msg.MessageText)

			// Regulate amount being a non number type error.
			if err != nil {
				app.TimeStampErrorPrint("Amount received by TransactionHandler is not a number.")
				return
			}

			// Add the bitcakes from the transaction.
			transactionalHandler.BitcakeManager.AddSomeBitcakes(amount)

			// Lock the color lock.
			app.ColorLock.Lock()
			// Release lock at the end of the scope.
			defer app.ColorLock.Unlock()
			// Record the amount got by the transaction.
			if bm, ok := transactionalHandler.BitcakeManager.(*LaiYangBitcakeManager); ok {
				if msg, ok1 := transactionalHandler.ClientMessage.(*message.BasicMessage); ok1 && msg.White {
					bm.RecordGetTransaction(msg.OriginalSenderInfo.Id, amount)
				}
			}
		}
	} else {
		app.TimeStampErrorPrint("Transaction handler got message of non BasicMessage type.")
	}

}
