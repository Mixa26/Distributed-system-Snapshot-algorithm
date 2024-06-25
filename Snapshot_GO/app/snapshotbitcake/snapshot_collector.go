package snapshotbitcake

type SnapshotCollector interface {
	Run()
	GetBitcakeManager() BitcakeManager
	AddLYSnapshotInfo(id int, lySnapshotResult *LYSnapshotResult)
	StartCollecting()
	Stop()
}
