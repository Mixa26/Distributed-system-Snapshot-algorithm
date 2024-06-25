package command

import (
	"Snapshot/app"
	"strconv"
	"strings"
	"time"
)

type Pause_command struct{}

func (cmd Pause_command) CommandName() string {
	return "pause"
}

func (cmd Pause_command) Execute(args string) {
	timeToSleep, err := strconv.Atoi(strings.TrimSpace(args))

	// Check if the arg is a number.
	if err != nil {
		app.TimeStampErrorPrint("FAILED TO PARSE NUMBER IN PAUSE COMMAND, PROVIDED: " + args)
		return
	}

	// Check that the number is positive
	if timeToSleep < 0 {
		app.TimeStampErrorPrint("TIME TO SLEEP ARG IN PAUSE COMMAND MUST BE POSITIVE! PROVIDED: " + args)
		return
	}

	app.TimeStampPrint("Pausing for " + strconv.Itoa(timeToSleep) + "ms.")

	// Sleep for timeToSleep miliseconds.
	time.Sleep(time.Duration(time.Duration(timeToSleep).Milliseconds()))
}
