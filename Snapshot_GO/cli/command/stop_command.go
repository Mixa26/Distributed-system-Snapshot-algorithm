package command

import (
	"Snapshot/app"
	"Snapshot/app/snapshotbitcake"
	"Snapshot/servent"
)

type Stop_command struct {
	parser            *CLIParser
	listener          *servent.SimpleServentListener
	snapshotCollector *snapshotbitcake.SnapshotCollector
}

func (cmd Stop_command) CommandName() string {
	return "stop"
}

func (cmd Stop_command) Execute(args string) {
	app.TimeStampPrint("Stopping...")
	cmd.parser.Stop()
	cmd.listener.Stop()
	// We have to cast the pointer to interface into a interface.
	if sc, ok := (*cmd.snapshotCollector).(snapshotbitcake.SnapshotCollector); ok {
		sc.Stop()
	}
}
