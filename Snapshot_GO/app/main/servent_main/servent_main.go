package main

import (
	"Snapshot/app"
	"Snapshot/app/snapshotbitcake"
	"Snapshot/cli/command"
	"Snapshot/servent"
	"os"
	"strconv"
	"sync"
)

var wgListener sync.WaitGroup

/*
This is going to be an individual process running a servent.
Command line arg at 0 is path to servent list file, 1 servent's id.
Since os.Args argument at pos:
0 - is the exe of the compiled file
1 - actual argument 0 that was passed by us
2 - actual argument 1 that was passed by us and so on...
*/
func main() {
	// The number is 3 as mentioned because 0 - exe, 1 - file for servent config, 2 - servent id.
	if len(os.Args) != 3 {
		app.ExitWithErrorMsg("Provide the servent list file along with the id of the servent!")
	}

	serventListFile := os.Args[1]

	// We read the config file and collect servents into a list in app_config.
	app.ReadConfig(serventListFile)

	// Read servents id from args.
	serventId, err := strconv.Atoi(os.Args[2])

	if err != nil {
		app.ExitWithErrorMsg("COULDN'T PARSE ID OF SERVENT! BAD ARGUMENT: " + os.Args[2])
	}

	if serventId < 0 || serventId > app.SERVENT_COUNT {
		app.ExitWithErrorMsg("INVALID SERVENT ID PROVIDED: " + os.Args[2])
	}

	// Set current process servent info in app_config.
	app.MyServentInfo = *app.GetServentInfo(serventId)

	if app.MyServentInfo.Port < 1000 || app.MyServentInfo.Port > 2000 {
		app.ExitWithErrorMsg("Port number should be in range 1000-2000.")
	}

	app.TimeStampPrint("Starting servent " + strconv.Itoa(app.MyServentInfo.Id))

	// Unlike other languages go doesn't wait for it's threads to finish
	// their tasks >:(
	// So we need to specify the main thread to wait for other threads.
	var wg sync.WaitGroup
	// We have to wait for 3 goroutines (threads).
	// Those are SnapshotCollector, SimpleServentListener and CLIParser.
	wg.Add(3)

	// Snapshot collector to be used in the system.
	var snapshotCollector snapshotbitcake.SnapshotCollector

	// Start snapshoCollector as thread.
	snapshotCollector = snapshotbitcake.ConstructSnapshotCollectorWorker()
	go func() {
		defer wg.Done()
		snapshotCollector.Run()
	}()

	wgListener.Add(1)

	// Run the simpleServentListener thread.
	simpleServentListener := servent.ConstructSimpleServentListener(snapshotCollector)
	go func() {
		defer wg.Done()
		simpleServentListener.Run(&wgListener)
	}()

	wgListener.Wait()

	// Construct CLI and run it as a thread.
	CLIParser := command.ConstructCLIParser(simpleServentListener, &snapshotCollector)
	go func() {
		defer wg.Done()
		CLIParser.Run()
	}()

	// Wait for the other threads to finish.
	wg.Wait()
}
