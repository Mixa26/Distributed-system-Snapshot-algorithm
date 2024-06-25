package app.snapshot_bitcake;

import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.BiFunction;

import app.AppConfig;
import servent.message.Message;
import servent.message.snapshot.LYMarkerMessage;
import servent.message.snapshot.LYTellMessage;
import servent.message.util.MessageUtil;

public class LaiYangBitcakeManager implements BitcakeManager {

	private final AtomicInteger currentAmount = new AtomicInteger(1000);

	public void takeSomeBitcakes(int amount) {
		currentAmount.getAndAdd(-amount);
	}

	public void addSomeBitcakes(int amount) {
		currentAmount.getAndAdd(amount);
	}

	public int getCurrentBitcakeAmount() {
		return currentAmount.get();
	}

	//private Map<Integer, Integer> giveHistory = new ConcurrentHashMap<>();
	//private Map<Integer, Integer> getHistory = new ConcurrentHashMap<>();

	// Since we're using the Li extension of the algorithm
	// we don't need the histories, but rather the state
	// of the canal (send bitcakes - received bitcakes)
	// for each canal connection with another servent.
	private final Map<Integer, Pair> transactions;

	private final Object MASTER_LOCK = new Object();

	private final AtomicInteger masterId = new AtomicInteger(-1);

	private final AtomicInteger parentId = new AtomicInteger(-1);

	private AtomicInteger markersIveReceived = new AtomicInteger(0);

	private LYSnapshotResult snapshotResult;

	private List<LYSnapshotResult> childrenSnapshotResults;

	private Set<Integer> idBorderSet;


	private final Object LOCK = new Object();

	public LaiYangBitcakeManager() {
		transactions = new HashMap<>();
		idBorderSet = new HashSet<>();
		childrenSnapshotResults = new ArrayList<>();
		clearHistory();
	}

	/*
	 * This value is protected by AppConfig.colorLock.
	 * Access it only if you have the blessing.
	 */
	public int recordedAmount = 0;

	public void markerEvent(Pair snapshotID, int senderId, SnapshotCollector snapshotCollector, LYTellMessage childInfo) {
		//synchronized (AppConfig.colorLock) {
		//AppConfig.isWhite.set(false);

		//Spezialetti-Kearns algorithm
		synchronized (MASTER_LOCK) {
			int masterIdCopy = masterId.get();
			if (masterIdCopy == -1) {
				// No one is my master so someone is my new master/I will be the master.
				masterId.set(snapshotID.first);
				masterIdCopy = snapshotID.first;
				if (masterIdCopy == AppConfig.myServentInfo.getId()) {
					markersIveReceived.decrementAndGet();
				}
				// Also set my parent to the servent i've received the marker from,
				// so i know whom to send the snapshot result to.
				parentId.set(senderId);
			} else if (masterIdCopy != snapshotID.first){
				// Initiator tried to initialize a snapshot, but a master has been already given to him.
				if (AppConfig.myServentInfo.getId() == snapshotID.first){
					System.out.println("CALLED RESET!!!!!!!!!!!!!!!!");
					// Hook to escape the snapshot and restart SnapshotCollectorWorker.
					snapshotCollector.snapshotDuringSnapshotExit();
					return;
				}
				// Another servent has contacted me from another region (he is on the border with me).
				// I need to remember his masters id so my master can contact him
				// to retrieve the snapshot information from that other region.
				System.out.println("RECEIVED FOREIGN BORDER ID: " + snapshotID.first);
				idBorderSet.add(snapshotID.first);
				// I already have a parent.
			}

			// Increment the number of markers so i know when the snapshot collection is over and
			// I can send my collected data to my parent.
			int collected = 0;

			collected = markersIveReceived.incrementAndGet();

			System.out.println("-----------------------------------------");
			System.out.println("MARKER RECEIVED FROM " + senderId);
			System.out.println("MY MASTER: " + masterIdCopy);
			System.out.println("MY PARENT: " + parentId.get());
			System.out.println("COLLECTED MARKERS: " + collected);
			System.out.println("-----------------------------------------");

			System.out.println("TEST1 " + (childInfo != null));
			System.out.println("TEST2 " + (masterIdCopy == snapshotID.first) + " MASTER IS " + masterIdCopy + " AND SNAPSHOTID.FIRST IS " + snapshotID.first);
			// Ensure we don't add info of servent from another region.
			if (childInfo != null && masterIdCopy == snapshotID.first){
				System.out.print("IM ADDING TO SNAPSHOT RES THIS: ");
				for (LYSnapshotResult result : childInfo.getLYSnapshotResult()) {
					System.out.println(result.getServentId());
				}
				// Add snapshot info for the neighbor if he is from my region.
				childrenSnapshotResults.addAll(childInfo.getLYSnapshotResult());
				// Add borders from possible child.
				addIdBorderSet(childInfo.getIdBorderSet());
			}

			// If i received all my neighbor markers the collection is over and i can send my snapshot data to my parent.
			if (collected == AppConfig.myServentInfo.getNeighbors().size()){

				// Don't send a message to myself if im the master, i've already collected all my info.
				if (masterIdCopy != AppConfig.myServentInfo.getId()) {
					// Add children's snapshots.
					List<LYSnapshotResult> snapshotResults = new ArrayList<>(childrenSnapshotResults);
					snapshotResults.add(snapshotResult);

					System.out.print("My children snapshots are: ");
					for (LYSnapshotResult result : snapshotResults){
						System.out.print(result.getServentId() + " ");
					}
					System.out.println();

					Set<Integer> idBorderSetCopy = new HashSet<>(idBorderSet);
					LYTellMessage tellMessage = new LYTellMessage(
							snapshotID.second, AppConfig.myServentInfo, AppConfig.getInfoById(snapshotID.first), snapshotResults, idBorderSetCopy);

					Message clMarker = new LYMarkerMessage(snapshotID.second, AppConfig.myServentInfo, AppConfig.getInfoById(parentId.get()), masterIdCopy, tellMessage);
					System.out.println("GOT ALL MARKERS SENDING RESULT TO PARENT: " + parentId.get());
					MessageUtil.sendMessage(clMarker);
					restartSpezialleti();
				} else {
					System.out.println("I've received all my children info!");
					System.out.println("Neighbor regions are: " + idBorderSet);
					System.out.print("My children snapshots are: ");
					for (LYSnapshotResult result : childrenSnapshotResults){
						System.out.print(result.getServentId() + " ");
					}
					System.out.println();
					snapshotCollector.addLYSnapshotInfos(snapshotID.second, childrenSnapshotResults, idBorderSet);
					restartSpezialleti();
					((SnapshotCollectorWorker)snapshotCollector).stopWaiting();
				}

				return;
			}
		}

		// Ensure we don't send multiple markers for one snapshot.
		if (AppConfig.initiatorsSnapshotNumber.get(String.valueOf(masterId)) < snapshotID.second) {
			// Increment the snapshot nth number for the initiator, so we know we received the marker for that snapshot.
			AppConfig.initiatorsSnapshotNumber.put(String.valueOf(masterId), snapshotID.second);
		} else {
			// If servents contact us to make a snapshot for ith snapshot, but we already did, return.
			// Also this makes sure we don't propagate a message from a neighbor region with another master.
			return;
		}

		Map<Integer, Pair> recordedTransactions = null;
		synchronized (LOCK) {
			recordedAmount = getCurrentBitcakeAmount();
			recordedTransactions = deepCopyTransactions(transactions);
			// Clear my channels.
			//clearHistory();
		}


		snapshotResult = new LYSnapshotResult(
				AppConfig.myServentInfo.getId(), recordedAmount, recordedTransactions);

        /*synchronized (LOCK) {
			System.out.println(recordedAmount);
			int sum = 0;
			for (Map.Entry<Integer, Pair> transaction : transactions.entrySet()){
				System.out.println(transaction.getKey() + " " + transaction.getValue().first + " " + transaction.getValue().second);
				sum += transaction.getValue().first - transaction.getValue().second;
			}
			System.out.println("total: " + (recordedAmount + sum));
		}*/

		if (snapshotID.first == AppConfig.myServentInfo.getId()) {
			snapshotCollector.addLYSnapshotInfos(
					snapshotID.second,
					Arrays.asList(snapshotResult),
					null);
		}
        /*else {
			Message tellMessage = new LYTellMessage(
					snapshotID.second, AppConfig.myServentInfo, AppConfig.getInfoById(snapshotID.first), snapshotResult);

			MessageUtil.sendMessage(tellMessage);
		}*/

		for (Integer neighbor : AppConfig.myServentInfo.getNeighbors()) {

			// Don't send a marker to my parent.
			if (neighbor == senderId) {
				continue;
			}

			Message clMarker = new LYMarkerMessage(snapshotID.second, AppConfig.myServentInfo, AppConfig.getInfoById(neighbor), snapshotID.first, null);
			MessageUtil.sendMessage(clMarker);
			try {
				/*
				 * This sleep is here to artificially produce some white node -> red node messages.
				 * Not actually recommended, as we are sleeping while we have colorLock.
				 */
				Thread.sleep(100);

			} catch (InterruptedException e) {
				e.printStackTrace();
			}
		}
		//}
	}

	public void restartSpezialleti() {
		// Since the collection is over clear the data for the next snapshot.
		masterId.set(-1);
		parentId.set(-1);
		markersIveReceived.set(0);
		idBorderSet.clear();
		snapshotResult = null;
		childrenSnapshotResults = new ArrayList<>();
	}

	private Map<Integer, Pair> deepCopyTransactions(Map<Integer, Pair> originalMap) {
		Map<Integer, Pair> copyMap = new HashMap<>();
		for (Map.Entry<Integer, Pair> entry : originalMap.entrySet()) {
			int key = entry.getKey();
			Pair value = entry.getValue();
			copyMap.put(key, new Pair(value.first, value.second));
		}
		return copyMap;
	}


	public void recordGiveTransaction(int neighbor, int amount) {
		//giveHistory.compute(neighbor, new MapValueUpdater(amount));
		synchronized (LOCK) {
			Pair pair = transactions.get(neighbor);
			pair.first += amount;
			transactions.put(neighbor, pair);
		}
	}

	public void recordGetTransaction(int neighbor, int amount) {
		//getHistory.compute(neighbor, new MapValueUpdater(amount));
		synchronized (LOCK) {
			Pair pair = transactions.get(neighbor);
			pair.second += amount;
			transactions.put(neighbor, pair);
		}
	}

	public void addIdBorderSet(Set<Integer> idBorderSet){
		this.idBorderSet.addAll(idBorderSet);
	}

	public void clearHistory() {
		for (Integer neighbor : AppConfig.myServentInfo.getNeighbors()) {
			transactions.put(neighbor, new Pair(0, 0));
		}
	}
}
