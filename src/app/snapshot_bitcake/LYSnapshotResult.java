package app.snapshot_bitcake;

import java.io.Serializable;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Snapshot result for servent with id serventId.
 * The amount of bitcakes on that servent is written in recordedAmount.
 * The channel messages are recorded in giveHistory and getHistory.
 * In Lai-Yang, the initiator has to reconcile the differences between
 * individual nodes, so we just let him know what we got and what we gave
 * and let him do the rest.
 * 
 * @author bmilojkovic
 *
 */
public class LYSnapshotResult implements Serializable {

	private static final long serialVersionUID = 8939516333227254439L;
	
	private final int serventId;
	private int recordedAmount;
	//private final Map<Integer, Integer> giveHistory;
	//private final Map<Integer, Integer> getHistory;

	// Since we're using the Li extension of the algorithm
	// we don't need the histories, but rather the state
	// of the canal (send bitcakes - received bitcakes)
	// for each canal connection with another servent.
	private final Map<Integer, Pair> transactions;

	
	public LYSnapshotResult(int serventId, int recordedAmount, Map<Integer, Pair> transactions) {
		this.serventId = serventId;
		this.recordedAmount = recordedAmount;
		//this.giveHistory = new ConcurrentHashMap<>(giveHistory);
		//this.getHistory = new ConcurrentHashMap<>(getHistory);
		this.transactions = new ConcurrentHashMap<>(transactions);
	}
	public int getServentId() {
		return serventId;
	}
	public int getRecordedAmount() {
		return recordedAmount;
	}

	public Map<Integer, Pair> getTransactions() {
		return transactions;
	}

	public void setRecordedAmount(int recordedAmount) {
		this.recordedAmount = recordedAmount;
	}

	public int getSenderId(){
		return this.serventId;
	}
}
