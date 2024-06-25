package main

import (
	"Snapshot/app"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// ServentCLI used for force stopping non-stoped processes.
type ServentCLI struct {
	processes []*exec.Cmd
}

// Separate cli thread used for stopping processes if necessary.
func (serventCLI ServentCLI) run() {
	var input string

	for {
		// Wait for user input.
		fmt.Scanln(&input)

		// And if the command is stop kill all running processes.
		if strings.ToLower(input) == "stop" {
			for _, process := range serventCLI.processes {
				if err := process.Process.Kill(); err != nil {
					// This error can happen if you're not running this as
					// an administrator(Windows) or as su(Linux).
					log.Fatal("failed to kill process: ", err)
				}
			}
			break
		}
	}
}

func startServentTest(testName string) {
	// All the processes will be stored here.
	var serventProcesses []*exec.Cmd

	// Is the path from this file to the file where
	// servents will write.
	var BASE_DIR = "../../" + testName

	app.ReadConfig(BASE_DIR + "/servent_list.properties")

	fmt.Println("Starting multiple servent runner. If servents do not finish on their own, type \"stop\" to finish them")

	serventCount := app.SERVENT_COUNT

	for i := 0; i < serventCount; i++ {
		process := exec.Command("go", "run", "servent_main/servent_main.go", BASE_DIR+"/servent_list.properties", strconv.Itoa(i))

		// Make output, error folders if they do not exist.
		// The numbers are for the linux permissions rwx.
		// Keep in mind the input folder needs to be provided so we
		// can have input for the servents.
		os.Mkdir(BASE_DIR+"/output", 0777)
		os.Mkdir(BASE_DIR+"/error", 0777)

		// Files for output, error and input that should be provided in the folder of the properties file as defined.
		// If they exist, they will be truncated, otherwise created.
		// Keep in mind the input files need to be provided so that
		// servents know what to do.
		output, err1 := os.Create(BASE_DIR + "/output/servent" + strconv.Itoa(i) + "_out.txt")
		err, err2 := os.Create(BASE_DIR + "/error/servent" + strconv.Itoa(i) + "_err.txt")
		input, err3 := os.Open(BASE_DIR + "/input/servent" + strconv.Itoa(i) + "_in.txt")

		// Check for errors with opening the files.
		if err1 != nil {
			app.ExitWithErrorMsg("BAD OUTPUT FILE: " + err1.Error())
		}

		if err2 != nil {
			app.ExitWithErrorMsg("BAD ERROR FILE: " + err2.Error())
		}

		if err3 != nil {
			app.ExitWithErrorMsg("BAD INPUT FILE: " + err3.Error())
		}

		// Redirect the communication of the process to these files.
		process.Stdout = output
		process.Stderr = err
		process.Stdin = input

		// Start the process.
		process.Start()
		// Add the new process.
		serventProcesses = append(serventProcesses, process)
	}

	// Start CLI thread and wait for user to input "stop".
	var serventCLI ServentCLI
	serventCLI.processes = serventProcesses

	// Wait for cli to finish.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		serventCLI.run()
	}()
	wg.Wait()

	for _, process := range serventProcesses {
		process.Wait()
	}

	app.TimeStampPrint("All servent processes finished.")
}

func main() {
	startServentTest("ly_snapshot")
}
