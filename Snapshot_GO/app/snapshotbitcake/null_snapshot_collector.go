package snapshotbitcake

type NullSnapshotCollector struct{}

func (nullSnapshotCollector NullSnapshotCollector) Run() {

}

func (nullSnapshotCollector NullSnapshotCollector) Stop() {

}

func (nullSnapshotCollector NullSnapshotCollector) GetBitcakeManager() BitcakeManager {
	return nil
}

func (nullSnapshotCollector NullSnapshotCollector) AddLYSnapshotInfo(id int, lySnapshotResult *LYSnapshotResult) {

}

func (nullSnapshotCollector NullSnapshotCollector) StartCollecting() {

}
