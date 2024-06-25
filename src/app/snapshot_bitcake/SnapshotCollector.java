package app.snapshot_bitcake;

import app.Cancellable;

import java.util.List;
import java.util.Set;

/**
 * Describes a snapshot collector. Made not-so-flexibly for readability.
 * 
 * @author bmilojkovic
 *
 */
public interface SnapshotCollector extends Runnable, Cancellable {

	BitcakeManager getBitcakeManager();

	void addLYSnapshotInfos(int snapshotNumber, List<LYSnapshotResult> lySnapshotResult, Set<Integer> idBorderSet);

	void snapshotDuringSnapshotExit();

	void startCollecting();

}