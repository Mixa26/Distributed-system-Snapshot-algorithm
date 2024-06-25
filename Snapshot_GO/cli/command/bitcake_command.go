package command

import (
	"Snapshot/app"
	"Snapshot/app/snapshotbitcake"
)

type Bitcake_command struct {
	collector snapshotbitcake.SnapshotCollector
}

func (cmd Bitcake_command) CommandName() string {
	return "bitcake_info"
}

func (cmd Bitcake_command) Execute(args string) {
	app.TimeStampPrint("Starting bitcake collecting...")
	cmd.collector.StartCollecting()
}
