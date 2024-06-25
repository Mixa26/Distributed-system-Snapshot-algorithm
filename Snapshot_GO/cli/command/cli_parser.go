package command

import (
	"Snapshot/app"
	"Snapshot/app/snapshotbitcake"
	"Snapshot/servent"
	"bufio"
	"os"
	"strings"
	"sync"
)

type CLIParser struct {
	commandList []CliCommand
	working     bool
}

// Constructor for CLIParser.
// We add all the existing commands to the CLI.
func ConstructCLIParser(
	simpleServentListener *servent.SimpleServentListener,
	snapshotCollector *snapshotbitcake.SnapshotCollector,
) *CLIParser {
	newCLIParser := CLIParser{}

	newCLIParser.working = true

	commandList := []CliCommand{
		Info_command{},
		Pause_command{},
		Transaction_burst_command{(*snapshotCollector).GetBitcakeManager()},
		Bitcake_command{*snapshotCollector},
		Stop_command{&newCLIParser, simpleServentListener, snapshotCollector},
	}

	newCLIParser.commandList = commandList

	return &newCLIParser
}

// We have to wait for the execution of threads.
// Thats what wg is for.
var wg sync.WaitGroup

// Function to run CLIParser as a thread.
func (cliParser CLIParser) Run() {
	var input string
	var command string
	var args string

	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Quit the run loop if working set to false by outside.
		if !cliParser.working {
			break
		}

		// Scan the input line from the file if these is any more.
		if scanner.Scan() {
			// Get the new input line.
			input = scanner.Text()

			indexOfSpace := strings.Index(input, " ")

			// Determine the command provided.
			if indexOfSpace != -1 {
				command = input[:indexOfSpace]
				args = input[indexOfSpace:]
			} else {
				command = input
			}

			found := false

			// Find the command in the commandList list.
			for _, cliCommand := range cliParser.commandList {
				if cliCommand.CommandName() == command {
					cliCommand.Execute(args)
					found = true
					break
				}
			}

			// If command is not found, inform the user.
			if !found {
				app.TimeStampErrorPrint("Command not found! Provided: " + command)
			}
		}
	}
	// Wait for possible transaction burst workers.
	wg.Wait()
}

func (cliParser *CLIParser) Stop() {
	(*cliParser).working = false
}
