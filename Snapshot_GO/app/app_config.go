package app

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Format for time stamp printing.
const timeFormat = "15:04:05"

// IP on which servents will listen
var SERVENT_IP string = "localhost"

// Number of servents in the system property.
var SERVENT_COUNT int

// Is the graph a clique property.
var IS_CLIQUE bool

// Type of snapshot to be used in the system.
var SnapshotType string

// Is the servent:
// White - Not in snapshot regime.
// Red - In snapshot regime.
var IsWhite atomic.Bool

// Lock for the servent color checking, the variable IsWhite.
var ColorLock *sync.Mutex

// Servent list read from properties.
var servents map[int]*ServentInfo

// Servent ids that can initiate a snapshot.
var SnapshotInitiatorServentIds []int

// Servent info of current process.
var MyServentInfo ServentInfo

/*
Prints a message with a time stamp at the beginning.
$ - normal print
@ - non fatal error
X - fatal error
*/
func TimeStampPrint(message string) {
	fmt.Println(" $ " + time.Now().Format(timeFormat) + " - " + message)
}

func TimeStampErrorPrint(message string) {
	log.Println(" @ " + time.Now().Format(timeFormat) + " - " + message)
}

func ExitWithErrorMsg(message string) {
	log.Println(" X " + time.Now().Format(timeFormat) + " - " + message + " | Exiting...")
	os.Exit(0)
}

/*
Reads config for setting up servents.
File is of format:
servent_count=2 /number of servets
clique=false /is the graph a clique?
fifo=false /is the message queue fifo?
servent0.port=1100 //ports and neighbors of servents
servent1.port=1200
servent0.neighbors=1
servent1.neighbors=0
*/
func ReadConfig(fileName string) {
	file, err := os.Open(fileName)

	// Error finding or reading the file.
	if err != nil {
		TimeStampErrorPrint("Make sure you're starting the app from the dir of multiple_servent_starter.go with go run multiple_servent_starter.go")
		ExitWithErrorMsg("ERROR READING FILE: " + fileName)
	}
	// Close the file at the end of this scope!
	defer file.Close()

	// Initialize the servent list.
	servents = make(map[int]*ServentInfo)
	// Initialize the lock.c
	ColorLock = &sync.Mutex{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineSlices := strings.Split(scanner.Text(), "=")

		if len(lineSlices) == 2 {
			switch lineSlices[0] {
			case "servent_count":
				SERVENT_COUNT, err = strconv.Atoi(lineSlices[1])

				// Servent count isn't defined properly (exp. servent_count=abc)
				if err != nil {
					ExitWithErrorMsg("FAILED PARSING SERVENT COUNT PROPERTY LINE: " + scanner.Text())
				}
			case "clique":
				IS_CLIQUE, err = strconv.ParseBool(lineSlices[1])

				// Clique value wasnt define properly (exp. clique=yes).
				if err != nil {
					ExitWithErrorMsg("FAILED PARSING CLIQUE PROPERTY AT LINE: " + scanner.Text())
				}
			case "initiators":
				initiatorStrs := strings.Split(lineSlices[1], ",")

				var initiatorInts []int

				// Parse all the initiators from line initiators=0,1,2 so we get 0,1,2.
				for _, initiator := range initiatorStrs {
					initiatorInt, err := strconv.Atoi(initiator)

					// A initiator is not a number (exp. initiators=0,a,2).
					if err != nil {
						ExitWithErrorMsg("FAILED PARSING INITIATORS AT LINE: " + scanner.Text())
					}

					initiatorInts = append(initiatorInts, initiatorInt)
				}

				// We set initiators only if all of them are numbers.
				SnapshotInitiatorServentIds = initiatorInts
			default:
				if strings.HasPrefix(lineSlices[0], "servent") {
					// Here we grap the port id (exp. servent0.port=1100/servent0.neighbors=1,2,6, we grab 0)
					splitByDot := strings.Split(lineSlices[0], ".")
					serventId, err := strconv.Atoi(strings.Split(splitByDot[0], "servent")[1])

					// Id is not a number (exp. serventb.port=1100)
					if err != nil {
						ExitWithErrorMsg("FAILED PARSING SERVENT ID PROPERTY LINE: " + scanner.Text())
					}

					// Try to parse the port.
					if splitByDot[1] == "port" {
						serventPort, err := strconv.Atoi(lineSlices[1])

						// Port is not a number (exp. serventb.port=1ca0)
						if err != nil {
							ExitWithErrorMsg("FAILED PARSING SERVENT PORT PROPERTY LINE: " + scanner.Text())
						}

						// Construct the new servent.
						newServentInfo := ConstructServentInfo(serventId, SERVENT_IP, serventPort, nil)

						servents[newServentInfo.Id] = newServentInfo
					} else if splitByDot[1] == "neighbors" || splitByDot[1] == "neighbours" {
						// Id of the servent to which we add neighbors.
						serventToAddNeighborsIndex := -1

						// Find the servent to which we want to add neighbors.
						for index, servent := range servents {
							if servent.Id == serventId {
								serventToAddNeighborsIndex = index
							}
						}

						// The neighbor line was defined before the port (exp. servent0.neighbors=1,2,6 then servent0.port=1100)
						if serventToAddNeighborsIndex == -1 {
							ExitWithErrorMsg("SERVENT NEIGHBORS MUST BE DEFINED AFTER SERVENT PORT: " + scanner.Text())
						}

						// Split the neighbors.
						neighborsStrs := strings.Split(lineSlices[1], ",")
						var neighborsInts []int

						// Go through all the neighbors.
						for _, neighbor := range neighborsStrs {
							num, err := strconv.Atoi(neighbor)

							// Neighbor is not a number (exp. servent0.neighbors=1,d,6).
							if err != nil {
								ExitWithErrorMsg("ERROR PARSING NEIGHBOR IN LINE: " + scanner.Text())
							}

							// Everything is cool, add the neighbor.
							neighborsInts = append(neighborsInts, num)
						}

						// We cannot change a part of the map element so we change it externally
						// and then assign the new version of the element.
						serventToChange := servents[serventToAddNeighborsIndex]
						serventToChange.Neighbors = neighborsInts
						servents[serventToAddNeighborsIndex] = serventToChange
					}
				} else {
					ExitWithErrorMsg("PROBLEM READING PROPERTIES, BAD LINE: " + scanner.Text())
				}
			}
		} else {
			ExitWithErrorMsg("PROBLEM READING PROPERTIES, BAD LINE: " + scanner.Text())
		}
	}
}

func GetServentInfo(id int) *ServentInfo {
	if id < 0 || id >= SERVENT_COUNT {
		TimeStampPrint("NON-EXISTENT SERVENT ID!")
	}
	return servents[id]
}
