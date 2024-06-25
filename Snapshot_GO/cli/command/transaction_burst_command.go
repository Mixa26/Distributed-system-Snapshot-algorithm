package command

import (
	"Snapshot/app"
	"Snapshot/app/snapshotbitcake"
	"math/rand"
	"time"
)

type Transaction_burst_command struct {
	bitcakeManager snapshotbitcake.BitcakeManager
}

var TRANSACTION_COUNT = 5
var BURST_WORKERS = 10
var MAX_TRANSFER_AMOUNT = 10

func (cmd Transaction_burst_command) run() {
	for i := 0; i < TRANSACTION_COUNT; i++ {
		for _, neighbor := range app.MyServentInfo.Neighbors {
			neighborInfo := app.GetServentInfo(neighbor)

			rand.New(rand.NewSource(time.Now().Unix()))

			amount := 1 + rand.Intn(MAX_TRANSFER_AMOUNT)

			/*
			* The message itself will reduce our bitcake count as it is being sent.
			* The sending might be delayed, so we want to make sure we do the
			* reducing at the right time, not earlier.
			 */
			transactionMessage := snapshotbitcake.ConstructTransactionMessage(&app.MyServentInfo, neighborInfo, amount, cmd.bitcakeManager)

			snapshotbitcake.SendMessage(transactionMessage)
		}
	}
}

func (cmd Transaction_burst_command) CommandName() string {
	return "transaction_burst"
}

func (cmd Transaction_burst_command) Execute(args string) {
	for i := 0; i < BURST_WORKERS; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cmd.run()
		}()
	}
}
