package snapshotbitcake

import (
	"Snapshot/app"
	"sync"
	"sync/atomic"
	"time"
)

type LaiYangBitcakeManager struct {
	CurrentAmount    atomic.Int64
	GiveHistoryLocks map[int]*sync.Mutex
	GiveHistory      map[int]int
	GetHistoryLocks  map[int]*sync.Mutex
	GetHistory       map[int]int
}

// Constructor for LaiYangBitcakeManager.
func ConstructLaiYankBitcakeManager() *LaiYangBitcakeManager {
	// Construct a new LaiYankBitcakeManager
	newLaiYankBitcakeManager := LaiYangBitcakeManager{atomic.Int64{}, make(map[int]*sync.Mutex), make(map[int]int), make(map[int]*sync.Mutex), make(map[int]int)}
	// Default starting value of bitcakes.
	newLaiYankBitcakeManager.CurrentAmount.Add(1000)
	// Initialize bitcake amounts for neighbors to 0.
	for _, neighbor := range app.MyServentInfo.Neighbors {
		newLaiYankBitcakeManager.GiveHistoryLocks[neighbor] = &sync.Mutex{}
		newLaiYankBitcakeManager.GiveHistory[neighbor] = 0
		newLaiYankBitcakeManager.GetHistoryLocks[neighbor] = &sync.Mutex{}
		newLaiYankBitcakeManager.GetHistory[neighbor] = 0
	}

	return &newLaiYankBitcakeManager
}

func (laiYangBitcakeManager *LaiYangBitcakeManager) TakeSomeBitcakes(amount int) {
	laiYangBitcakeManager.CurrentAmount.Add(-int64(amount))
}

func (laiYangBitcakeManager *LaiYangBitcakeManager) AddSomeBitcakes(amount int) {
	laiYangBitcakeManager.CurrentAmount.Add(int64(amount))
}

func (laiYangBitcakeManager *LaiYangBitcakeManager) GetCurrentBitcakeAmount() int {
	return int(laiYangBitcakeManager.CurrentAmount.Load())
}

var recordedAmount = 0

func (laiYangBitcakeManager *LaiYangBitcakeManager) MarkerEvent(collectorId int, snapshotCollector SnapshotCollector) {
	// Acquire the lock and make sure to unlock it at the end of this scope.
	app.ColorLock.Lock()
	defer app.ColorLock.Unlock()

	recordedAmount = laiYangBitcakeManager.GetCurrentBitcakeAmount()

	// Set servent color to red.
	app.IsWhite.CompareAndSwap(true, false)

	// Snapshot results.
	snapshotResult := LYSnapshotResult{
		ServentId:      app.MyServentInfo.Id,
		RecorderAmount: laiYangBitcakeManager.GetCurrentBitcakeAmount(),
		GiveHistory:    laiYangBitcakeManager.GiveHistory,
		GetHistory:     laiYangBitcakeManager.GetHistory,
	}

	// Add the snapshot info to the snapshot collector
	// if the collector is for my id.
	if collectorId == app.MyServentInfo.Id {
		snapshotCollector.AddLYSnapshotInfo(
			app.MyServentInfo.Id,
			&snapshotResult,
		)
	} else {
		SendMessage(ConstructLYTellMessage(
			&app.MyServentInfo,
			app.GetServentInfo(collectorId),
			&snapshotResult,
		))
	}

	// This is a thread wait group so we can wait
	// for execution of threads on line 94.
	var wg sync.WaitGroup

	for _, neighbor := range app.MyServentInfo.Neighbors {
		// Message my neighbors about the snapshot collection.
		clMarker := ConstructLYMarkerMessage(&app.MyServentInfo, app.GetServentInfo(neighbor), collectorId)
		// We have to wait for the threads to finish.
		wg.Add(1)
		go func() {
			defer wg.Done()
			SendMessage(&clMarker)
		}()
		// This sleep is here to artificially produce some white node -> red node messages.
		time.Sleep(time.Duration(time.Duration(100).Milliseconds()))
	}
	// Wait for all threads from sendMessage() calls above.
	wg.Wait()
}

func (laiYangBitcakeManager *LaiYangBitcakeManager) RecordGiveTransaction(neighbor int, amount int) {
	// We lock the map element only, not the whole map.
	laiYangBitcakeManager.GiveHistoryLocks[neighbor].Lock()
	// Unlock the element at the end of this scope.
	defer laiYangBitcakeManager.GiveHistoryLocks[neighbor].Unlock()
	// Update the value.
	laiYangBitcakeManager.GiveHistory[neighbor] += amount
}

func (laiYangBitcakeManager *LaiYangBitcakeManager) RecordGetTransaction(neighbor int, amount int) {
	// We lock the map element only, not the whole map.
	laiYangBitcakeManager.GetHistoryLocks[neighbor].Lock()
	// Unlock the element at the end of this scope.
	defer laiYangBitcakeManager.GetHistoryLocks[neighbor].Unlock()
	// Update the value.
	laiYangBitcakeManager.GetHistory[neighbor] += amount
}
