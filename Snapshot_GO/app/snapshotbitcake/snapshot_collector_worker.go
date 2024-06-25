package snapshotbitcake

import (
	"Snapshot/app"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type SnapshotCollectorWorker struct {
	working               bool
	Collecting            atomic.Bool
	CollectedLYValuesLock map[int]*sync.Mutex
	CollectedLYValues     map[int]LYSnapshotResult
	BitcakeManager        BitcakeManager
}

func ConstructSnapshotCollectorWorker() *SnapshotCollectorWorker {
	newSnapshotCollectorWorker := SnapshotCollectorWorker{
		true,
		atomic.Bool{},
		make(map[int]*sync.Mutex),
		make(map[int]LYSnapshotResult),
		ConstructLaiYankBitcakeManager(),
	}

	// Initialize locks for synchronization when adding a new LYSnapshotResult.
	for i := 0; i < app.SERVENT_COUNT; i++ {
		newSnapshotCollectorWorker.CollectedLYValuesLock[i] = &sync.Mutex{}
	}

	newSnapshotCollectorWorker.Collecting.Store(false)

	return &newSnapshotCollectorWorker
}

func (snapshotCollectorWorker *SnapshotCollectorWorker) GetBitcakeManager() BitcakeManager {
	return snapshotCollectorWorker.BitcakeManager
}

// Userd in Run().
func ContainsInt(slice []int, value int) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func (snapshotCollectorWorker SnapshotCollectorWorker) Run() {
	for {
		// If someone called Stop(), jump out of the for loop.
		if !snapshotCollectorWorker.working {
			break
		}

		// Not collecting yet - just sleep until we start actual work, or finish.
		for {
			if snapshotCollectorWorker.Collecting.Load() {
				break
			}
			time.Sleep(1000)
			if !snapshotCollectorWorker.working {
				return
			}
		}

		/*
		* Collecting is done in three stages:
		* 1. Send messages asking for values
		* 2. Wait for all the responses
		* 3. Print result
		 */

		// 1. Send asks.
		if lybm, ok := snapshotCollectorWorker.BitcakeManager.(*LaiYangBitcakeManager); ok {
			lybm.MarkerEvent(app.MyServentInfo.Id, &snapshotCollectorWorker)
		}

		// 2. Wait for responses or finish.
		waiting := true
		for {
			if !waiting {
				break
			}

			fmt.Println("Stuck waiting")
			if len(snapshotCollectorWorker.CollectedLYValues) == app.SERVENT_COUNT {
				waiting = false
			}
			time.Sleep(1000)

			// If someone called Stop() terminate.
			if !snapshotCollectorWorker.working {
				return
			}
		}

		// 3. Print results.
		sum := 0
		// Lock every node result.
		for _, lock := range snapshotCollectorWorker.CollectedLYValuesLock {
			lock.Lock()
		}

		// Print recorder amounts on servents.
		for id, value := range snapshotCollectorWorker.CollectedLYValues {
			sum += value.RecorderAmount
			app.TimeStampPrint("Recorded bitcake amount for node id: " + strconv.Itoa(id) + " is " + strconv.Itoa(value.RecorderAmount))
		}

		// Print unreceived amounts from servents if there are any.
		for i := 0; i < app.SERVENT_COUNT; i++ {
			for j := 0; j < app.SERVENT_COUNT; j++ {
				if i != j {
					// Yes i know chaos implementation, dont judge...
					ijAmount := 0
					jiAmount := 0
					if ContainsInt(app.GetServentInfo(i).Neighbors, j) && ContainsInt(app.GetServentInfo(j).Neighbors, i) {
						ijAmount = snapshotCollectorWorker.CollectedLYValues[i].GiveHistory[j]
						jiAmount = snapshotCollectorWorker.CollectedLYValues[j].GetHistory[i]
					}

					if ijAmount != jiAmount {
						app.TimeStampPrint(fmt.Sprintf("Unreceived bitcake amount %d from servent %d to servent %d!", ijAmount-jiAmount, i, j))
						sum += ijAmount - jiAmount
					}
				}
			}
		}

		// Unlock every node result.
		for _, lock := range snapshotCollectorWorker.CollectedLYValuesLock {
			lock.Unlock()
		}
	}
}

func (snapshotCollectorWorker *SnapshotCollectorWorker) AddLYSnapshotInfo(id int, lySnapshotResult *LYSnapshotResult) {
	// Make sure to lock and lock the corresponding value.
	snapshotCollectorWorker.CollectedLYValuesLock[id].Lock()
	defer snapshotCollectorWorker.CollectedLYValuesLock[id].Unlock()
	// Set the value.
	snapshotCollectorWorker.CollectedLYValues[id] = *lySnapshotResult
}

func (snapshotCollectorWorker *SnapshotCollectorWorker) StartCollecting() {
	// Mark the start of collecting.
	swapped := snapshotCollectorWorker.Collecting.CompareAndSwap(false, true)

	// If the value was already true, means we interrupted a recording in progress.
	// In other words the swapped = false.
	if !swapped {
		app.TimeStampErrorPrint("Tried to start collecting before finished with the previous collection.")
	}
}

func (snapshotCollectorWorker *SnapshotCollectorWorker) Stop() {
	snapshotCollectorWorker.working = false
}
