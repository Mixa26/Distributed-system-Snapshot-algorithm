package app.snapshot_bitcake;

import java.util.*;
import java.util.Map.Entry;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;

import app.AppConfig;
import servent.message.Message;
import servent.message.snapshot.RoundMessage;
import servent.message.util.MessageUtil;

/**
 * Main snapshot collector class. Has support for Naive, Chandy-Lamport
 * and Lai-Yang snapshot algorithms.
 * 
 * @author bmilojkovic
 *
 */
public class SnapshotCollectorWorker implements SnapshotCollector {

	private volatile boolean working = true;

	private AtomicBoolean collecting = new AtomicBoolean(false);

	private final List<Map<Integer, LYSnapshotResult>> collectedLYValues;

	private Set<Integer> idBorderSet;

	private static AtomicInteger snapshotNumber = new AtomicInteger(0);

	private static final Object LOCK = new Object();

	private static AtomicBoolean TRIED_SNAPSHOT_DURING_SNAPSHOT = new AtomicBoolean(false);

	private static AtomicBoolean waiting = new AtomicBoolean(true);

	private BitcakeManager bitcakeManager;

	public SnapshotCollectorWorker() {
		bitcakeManager = new LaiYangBitcakeManager();
		collectedLYValues = new ArrayList<>();
		idBorderSet = new HashSet<>();
	}

	@Override
	public BitcakeManager getBitcakeManager() {
		return bitcakeManager;
	}

	@Override
	public void run() {
		while(working) {
			/*
			 * Not collecting yet - just sleep until we start actual work, or finish
			 */
			TRIED_SNAPSHOT_DURING_SNAPSHOT.set(false);
			while (collecting.get() == false) {
				try {
					Thread.sleep(1000);
				} catch (InterruptedException e) {
					// TODO Auto-generated catch block
					e.printStackTrace();
				}

				if (working == false) {
					return;
				}
			}

			/*
			 * Collecting is done in three stages:
			 * 1. Send messages asking for values
			 * 2. Wait for all the responses
			 * 3. Print result
			 */

			//1 send asks
			// Snapshot ID is composed of my id, and the nth number of the snapshot we're doing for my servent id.
			Pair snapshotID;
			snapshotID = new Pair(AppConfig.myServentInfo.getId(), snapshotNumber.getAndIncrement());
			synchronized (LOCK) {
				if (collectedLYValues.size() != snapshotID.second) {
					System.err.println("Snapshot number mismatch. This should never happen.");
				}
				collectedLYValues.add(new HashMap<>());
			}
			AppConfig.timestampedStandardPrint("STARTING SNAPSHOT FOR " + snapshotID);
			((LaiYangBitcakeManager)bitcakeManager).markerEvent(snapshotID, AppConfig.myServentInfo.getId(), this, null);

			if (TRIED_SNAPSHOT_DURING_SNAPSHOT.get()) {
				continue;
			}

			//2 wait for responses or finish
			while (waiting.get()) {

				try {
					Thread.sleep(1000);
				} catch (InterruptedException e) {
					e.printStackTrace();
				}

				if (working == false) {
					return;
				}
			}
			waiting.set(true);

			// 2.1 Inform other region masters, and receive their info

			synchronized (LOCK){
				Iterator<Integer> borderIterator = idBorderSet.iterator();
				while (borderIterator.hasNext()){
					Message toSend = new RoundMessage(AppConfig.myServentInfo, AppConfig.getInfoById(borderIterator.next()), collectedLYValues.get(snapshotNumber.get()-1));

					MessageUtil.sendMessage(toSend);
				}
			}

			while (true) {
				int collected = 0;
				synchronized (LOCK){
					collected = collectedLYValues.get(snapshotNumber.get()-1).size();
				}
				if (collected == AppConfig.getServentCount()) {
					break;
				}

				try {
					Thread.sleep(1000);
				} catch (InterruptedException e) {
					e.printStackTrace();
				}

				if (working == false) {
					return;
				}
			}

			//print
			int sum;
			sum = 0;
			Map<Integer, LYSnapshotResult> snapshotResult = null;
			Map<Integer, LYSnapshotResult> previousSnapshotResult = null;
			synchronized (LOCK) {
				idBorderSet.clear();
				snapshotResult = deepCopyRes(collectedLYValues.get(snapshotID.second));
				if (snapshotID.second > 0) {
					previousSnapshotResult = deepCopyRes(collectedLYValues.get(snapshotID.second - 1));
				}
			}

			for (Entry<Integer, LYSnapshotResult> nodeResult : snapshotResult.entrySet()) {
				// Get the last snapshot
				sum += nodeResult.getValue().getRecordedAmount();
				AppConfig.timestampedStandardPrint(
						"Recorded bitcake amount for " + nodeResult.getKey() + " = " + nodeResult.getValue().getRecordedAmount());
			}

            for(int i = 0; i < AppConfig.getServentCount(); i++) {
				for (int j = 0; j < AppConfig.getServentCount(); j++) {
					if (i != j) {
						if (AppConfig.getInfoById(i).getNeighbors().contains(j) &&
								AppConfig.getInfoById(j).getNeighbors().contains(i)) {
							int ijAmount = snapshotResult.get(i).getTransactions().get(j).first;
							int jiAmount = snapshotResult.get(i).getTransactions().get(j).second;
							int lastHistory = 0;
							if (snapshotID.second > 0) {
								lastHistory = previousSnapshotResult.get(i).getTransactions().get(j).first-
										previousSnapshotResult.get(i).getTransactions().get(j).second;
							}

							if (ijAmount != jiAmount) {
								String outputString = String.format(
										"Unreceived bitcake amount: %d from servent %d to servent %d",
										ijAmount - jiAmount, i, j);
								AppConfig.timestampedStandardPrint(outputString);
								sum += ijAmount - jiAmount;
							}
						}
					}
				}
			}

			AppConfig.timestampedStandardPrint("System bitcake count: " + sum);
			collecting.set(false);
		}
	}

	private Map<Integer, LYSnapshotResult> deepCopyRes(Map<Integer, LYSnapshotResult> originalMap) {
		Map<Integer, LYSnapshotResult> copyMap = new HashMap<>();
		for (Map.Entry<Integer, LYSnapshotResult> entry : originalMap.entrySet()) {
			int key = entry.getKey();
			LYSnapshotResult value = entry.getValue();
			copyMap.put(key, new LYSnapshotResult(value.getSenderId(), value.getRecordedAmount(), value.getTransactions()));
		}
		return copyMap;
	}

	public void stopWaiting() {
		waiting.set(false);
	}

	@Override
	public void snapshotDuringSnapshotExit() {
		TRIED_SNAPSHOT_DURING_SNAPSHOT.set(true);
		snapshotNumber.decrementAndGet();
		synchronized (LOCK) {
			collectedLYValues.remove(collectedLYValues.size() - 1);
		}
		collecting.set(false);
	}

	@Override
	public void addLYSnapshotInfos(int snapshotNumber, List<LYSnapshotResult> lySnapshotResult, Set<Integer> idBorderSet) {
		synchronized (LOCK) {
			for (LYSnapshotResult result : lySnapshotResult) {
				collectedLYValues.get(snapshotNumber).put(result.getServentId(), result);
			}
			if (idBorderSet != null) {
				this.idBorderSet = new HashSet<>(idBorderSet);
			}
		}
	}

	public void addRoundResult(Map<Integer, LYSnapshotResult> results) {
		int snapNum = snapshotNumber.get() - 1;
		synchronized (LOCK) {
			for (Entry<Integer, LYSnapshotResult> resultEntry : results.entrySet()) {
				collectedLYValues.get(snapNum).put(resultEntry.getKey(), resultEntry.getValue());
			}
		}
	}

	@Override
	public void startCollecting() {
		boolean oldValue = this.collecting.get();

		if (oldValue == true) {
			AppConfig.timestampedErrorPrint("Tried to start collecting before finished with previous.");
		} else{
			this.collecting.set(true);
		}
	}

	@Override
	public void stop() {
		working = false;
	}

}
